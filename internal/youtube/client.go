// Package youtube wraps the YouTube Data API v3 for trend + niche features.
// Quota: 1 unit per call (videos.list, channels.list), 100 units per search.
// Free tier is 10,000 units/day, so each ttb run costs <150 units.
package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const BaseURL = "https://www.googleapis.com/youtube/v3"

// Client is a tiny YouTube Data API v3 wrapper.
type Client struct {
	APIKey string
	HTTP   *http.Client
}

// New returns a Client. apiKey must be a Google Cloud project key with the
// YouTube Data API v3 enabled.
func New(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTP:   &http.Client{Timeout: 30 * time.Second},
	}
}

// Video is the public subset of a YouTube video resource that ttb cares about.
type Video struct {
	ID           string
	Title        string
	Description  string
	ChannelID    string
	ChannelTitle string
	PublishedAt  string
	CategoryID   string
	Tags         []string
	Duration     string // ISO 8601 (e.g. PT4M13S)
	ViewCount    int64
	LikeCount    int64
	CommentCount int64
}

// Channel is the public subset of a YouTube channel resource that ttb uses.
type Channel struct {
	ID              string
	Title           string
	Description     string
	CustomURL       string
	Country         string
	PublishedAt     string
	SubscriberCount int64
	ViewCount       int64
	VideoCount      int64
}

// Trending returns the most popular videos for a region.
// regionCode is ISO 3166-1 alpha-2 (e.g. "US", "ID", "JP").
// categoryID is optional (empty = all categories).
func (c *Client) Trending(ctx context.Context, regionCode, categoryID string, maxResults int) ([]Video, error) {
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 25
	}
	params := url.Values{
		"part":       {"snippet,statistics,contentDetails"},
		"chart":      {"mostPopular"},
		"regionCode": {regionCode},
		"maxResults": {fmt.Sprintf("%d", maxResults)},
		"key":        {c.APIKey},
	}
	if categoryID != "" {
		params.Set("videoCategoryId", categoryID)
	}
	var resp trendingResponse
	if err := c.get(ctx, "/videos", params, &resp); err != nil {
		return nil, err
	}
	out := make([]Video, len(resp.Items))
	for i, it := range resp.Items {
		out[i] = it.toVideo()
	}
	return out, nil
}

// ChannelByHandle looks up a channel by its @handle (e.g. "@mkbhd").
func (c *Client) ChannelByHandle(ctx context.Context, handle string) (*Channel, error) {
	handle = strings.TrimPrefix(handle, "@")
	if handle == "" {
		return nil, fmt.Errorf("youtube: empty handle")
	}
	params := url.Values{
		"part":      {"snippet,statistics"},
		"forHandle": {"@" + handle},
		"key":       {c.APIKey},
	}
	return c.fetchChannel(ctx, params)
}

// ChannelByID looks up a channel by its channel ID (UC...).
func (c *Client) ChannelByID(ctx context.Context, id string) (*Channel, error) {
	params := url.Values{
		"part": {"snippet,statistics"},
		"id":   {id},
		"key":  {c.APIKey},
	}
	return c.fetchChannel(ctx, params)
}

// ChannelVideos returns the latest videos from a channel, with statistics.
// Costs 2 quota units: 1 for search, 1 for the videos.list batch.
func (c *Client) ChannelVideos(ctx context.Context, channelID string, maxResults int) ([]Video, error) {
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 25
	}
	// 1. search for videos in channel
	searchParams := url.Values{
		"part":       {"snippet"},
		"channelId":  {channelID},
		"order":      {"date"},
		"type":       {"video"},
		"maxResults": {fmt.Sprintf("%d", maxResults)},
		"key":        {c.APIKey},
	}
	var searchResp searchResponse
	if err := c.get(ctx, "/search", searchParams, &searchResp); err != nil {
		return nil, fmt.Errorf("youtube: search: %w", err)
	}
	if len(searchResp.Items) == 0 {
		return nil, nil
	}
	// 2. batch fetch statistics for those video IDs
	ids := make([]string, 0, len(searchResp.Items))
	meta := make(map[string]searchItemMeta, len(searchResp.Items))
	for _, it := range searchResp.Items {
		ids = append(ids, it.ID.VideoID)
		meta[it.ID.VideoID] = searchItemMeta{
			Title:        it.Snippet.Title,
			Description:  it.Snippet.Description,
			ChannelID:    it.Snippet.ChannelID,
			ChannelTitle: it.Snippet.ChannelTitle,
			PublishedAt:  it.Snippet.PublishedAt,
		}
	}
	vidParams := url.Values{
		"part":       {"statistics,contentDetails"},
		"id":         {strings.Join(ids, ",")},
		"maxResults": {fmt.Sprintf("%d", maxResults)},
		"key":        {c.APIKey},
	}
	var vidResp trendingResponse
	if err := c.get(ctx, "/videos", vidParams, &vidResp); err != nil {
		return nil, fmt.Errorf("youtube: videos.batch: %w", err)
	}
	out := make([]Video, 0, len(vidResp.Items))
	for _, it := range vidResp.Items {
		m, ok := meta[it.ID]
		if !ok {
			continue
		}
		v := it.toVideo()
		v.Title = m.Title
		v.Description = m.Description
		v.ChannelID = m.ChannelID
		v.ChannelTitle = m.ChannelTitle
		v.PublishedAt = m.PublishedAt
		out = append(out, v)
	}
	return out, nil
}

// SearchVideos does a generic video search. Useful for future features.
func (c *Client) SearchVideos(ctx context.Context, query string, maxResults int) ([]Video, error) {
	if maxResults <= 0 || maxResults > 50 {
		maxResults = 25
	}
	params := url.Values{
		"part":       {"snippet"},
		"q":          {query},
		"type":       {"video"},
		"maxResults": {fmt.Sprintf("%d", maxResults)},
		"key":        {c.APIKey},
	}
	var resp searchResponse
	if err := c.get(ctx, "/search", params, &resp); err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(resp.Items))
	meta := make(map[string]searchItemMeta, len(resp.Items))
	for _, it := range resp.Items {
		ids = append(ids, it.ID.VideoID)
		meta[it.ID.VideoID] = searchItemMeta{
			Title:        it.Snippet.Title,
			Description:  it.Snippet.Description,
			ChannelID:    it.Snippet.ChannelID,
			ChannelTitle: it.Snippet.ChannelTitle,
			PublishedAt:  it.Snippet.PublishedAt,
		}
	}
	if len(ids) == 0 {
		return nil, nil
	}
	vidParams := url.Values{
		"part": {"statistics,contentDetails"},
		"id":   {strings.Join(ids, ",")},
		"key":  {c.APIKey},
	}
	var vidResp trendingResponse
	if err := c.get(ctx, "/videos", vidParams, &vidResp); err != nil {
		return nil, err
	}
	out := make([]Video, 0, len(vidResp.Items))
	for _, it := range vidResp.Items {
		m, ok := meta[it.ID]
		if !ok {
			continue
		}
		v := it.toVideo()
		v.Title = m.Title
		v.Description = m.Description
		v.ChannelID = m.ChannelID
		v.ChannelTitle = m.ChannelTitle
		v.PublishedAt = m.PublishedAt
		out = append(out, v)
	}
	return out, nil
}

// -- internal helpers --

func (c *Client) fetchChannel(ctx context.Context, params url.Values) (*Channel, error) {
	var resp channelResponse
	if err := c.get(ctx, "/channels", params, &resp); err != nil {
		return nil, err
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("youtube: channel not found")
	}
	it := resp.Items[0]
	return &Channel{
		ID:              it.ID,
		Title:           it.Snippet.Title,
		Description:     it.Snippet.Description,
		CustomURL:       it.Snippet.CustomURL,
		Country:         it.Snippet.Country,
		PublishedAt:     it.Snippet.PublishedAt,
		SubscriberCount: parseInt64(it.Statistics.SubscriberCount),
		ViewCount:       parseInt64(it.Statistics.ViewCount),
		VideoCount:      parseInt64(it.Statistics.VideoCount),
	}, nil
}

func (c *Client) get(ctx context.Context, path string, params url.Values, out any) error {
	u := BaseURL + path + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("youtube: new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("youtube: do: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("youtube: HTTP %d: %s", resp.StatusCode, trimForLog(string(body), 300))
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("youtube: decode: %w (body=%s)", err, trimForLog(string(body), 300))
	}
	return nil
}

func parseInt64(s string) int64 {
	var n int64
	_, _ = fmt.Sscanf(s, "%d", &n)
	return n
}

func trimForLog(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// -- response shapes (private, internal use only) --

type trendingItem struct {
	ID      string `json:"id"`
	Snippet struct {
		Title        string   `json:"title"`
		Description  string   `json:"description"`
		ChannelID    string   `json:"channelId"`
		ChannelTitle string   `json:"channelTitle"`
		PublishedAt  string   `json:"publishedAt"`
		CategoryID   string   `json:"categoryId"`
		Tags         []string `json:"tags"`
	} `json:"snippet"`
	Statistics struct {
		ViewCount    string `json:"viewCount"`
		LikeCount    string `json:"likeCount"`
		CommentCount string `json:"commentCount"`
	} `json:"statistics"`
	ContentDetails struct {
		Duration string `json:"duration"`
	} `json:"contentDetails"`
}

func (t trendingItem) toVideo() Video {
	return Video{
		ID:           t.ID,
		Title:        t.Snippet.Title,
		Description:  t.Snippet.Description,
		ChannelID:    t.Snippet.ChannelID,
		ChannelTitle: t.Snippet.ChannelTitle,
		PublishedAt:  t.Snippet.PublishedAt,
		CategoryID:   t.Snippet.CategoryID,
		Tags:         t.Snippet.Tags,
		Duration:     t.ContentDetails.Duration,
		ViewCount:    parseInt64(t.Statistics.ViewCount),
		LikeCount:    parseInt64(t.Statistics.LikeCount),
		CommentCount: parseInt64(t.Statistics.CommentCount),
	}
}

type trendingResponse struct {
	Items []trendingItem `json:"items"`
}

type channelResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			CustomURL   string `json:"customUrl"`
			Country     string `json:"country"`
			PublishedAt string `json:"publishedAt"`
		} `json:"snippet"`
		Statistics struct {
			ViewCount       string `json:"viewCount"`
			SubscriberCount string `json:"subscriberCount"`
			VideoCount      string `json:"videoCount"`
		} `json:"statistics"`
	} `json:"items"`
}

type searchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			Description  string `json:"description"`
			ChannelID    string `json:"channelId"`
			ChannelTitle string `json:"channelTitle"`
			PublishedAt  string `json:"publishedAt"`
		} `json:"snippet"`
	} `json:"items"`
}

type searchItemMeta struct {
	Title, Description, ChannelID, ChannelTitle, PublishedAt string
}
