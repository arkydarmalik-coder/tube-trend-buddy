// Package prompts holds the 8 system+user prompt templates used by ttb.
// One pair of functions per feature. Pure functions, no I/O.
package prompts

import "fmt"

// Title returns (system, user) prompts for the title generator.
func Title(niche, audience, lang string, count int) (string, string) {
	sys := "You are a YouTube growth strategist. You write click-worthy, " +
		"non-clickbait titles that balance curiosity and clarity. Avoid misleading claims. " +
		"Reply in the user's requested language. No preamble, no numbering, no markdown - " +
		"just one title per line, prefixed with a dash and a space."
	user := fmt.Sprintf(
		"Niche: %q\nTarget audience: %s\nLanguage: %s\nGenerate %d title ideas for videos in this niche. "+
			"Each title must be under 70 characters and trigger curiosity without misleading the viewer. "+
			"Mix proven patterns: how-to, listicle, question, contrarian, secret/reveal, story-driven.",
		niche, audience, lang, count)
	return sys, user
}

// Tags returns (system, user) prompts for the SEO tag optimizer.
func Tags(title, description, lang string, count int) (string, string) {
	sys := "You are a YouTube SEO specialist. You produce a balanced mix of " +
		"high-volume short tags, mid-tail tags, and specific long-tail tags. " +
		"Reply in the user's language. Output one tag per line, prefixed with a dash and a space. " +
		"No explanations, no numbering, no markdown fences."
	descLine := ""
	if description != "" {
		descLine = fmt.Sprintf("\nDescription: %s", description)
	}
	user := fmt.Sprintf(
		"Video title: %q%s\nLanguage: %s\nGenerate %d tags that maximize discoverability. "+
			"Sort loosely by likely search volume (broad first). Avoid # at the start - just the raw tag text.",
		title, descLine, lang, count)
	return sys, user
}

// Trend returns (system, user) prompts for the trend detector.
func Trend(region, category, period, lang string, count int) (string, string) {
	sys := "You are a YouTube trend analyst. You identify rising niches, " +
		"emerging formats, and untapped topic clusters. Reply in the user's language. " +
		"Output one trend per line: 'Topic -- why it is rising -- 1 suggested angle'. " +
		"No preamble, no markdown fences."
	user := fmt.Sprintf(
		"Region: %s\nCategory hint: %s\nLookback period: %s\nLanguage: %s\nList %d rising niches/topics. "+
			"For each, give a one-line rationale and a one-line suggested content angle that a mid-size creator could realistically win.",
		region, category, period, lang, count)
	return sys, user
}

// Niche returns (system, user) prompts for the channel/niche analyzer.
func Niche(channel string, deep bool, lang string) (string, string) {
	sys := "You are a YouTube strategy consultant. You give actionable, " +
		"specific advice backed by channel-format theory. Reply in the user's language. " +
		"Structure your reply in clear sections: 1) positioning summary, " +
		"2) content strengths, 3) content gaps, 4) 3 adjacent niches to expand into, " +
		"5) 90-day experiment ideas. Use plain text with section headers like 'POSITIONING:'."
	depthLine := "Standard analysis."
	if deep {
		depthLine = "Deep analysis - spend more time, include retention/CTA/CTR heuristics."
	}
	user := fmt.Sprintf(
		"Channel: %s\nDepth: %s\nLanguage: %s\nAudit this channel. If you do not know the exact channel, "+
			"infer from the handle and call out your assumption explicitly.",
		channel, depthLine, lang)
	return sys, user
}

// Describe returns (system, user) prompts for the description+hashtag writer.
func Describe(title, summary, timestamps string, includeCTA bool, lang string) (string, string) {
	sys := "You are a YouTube description copywriter. You write descriptions " +
		"that rank in search, surface in suggested, and drive action. Reply in the user's language. " +
		"Output plain text, no markdown fences, no preamble."
	tsLine := ""
	if timestamps != "" {
		tsLine = fmt.Sprintf("\nChapters (already in HH:MM format, separated by comma): %s", timestamps)
	}
	ctaLine := "Include a soft call-to-action (subscribe + next video)."
	if !includeCTA {
		ctaLine = "Do NOT include any call-to-action."
	}
	sumLine := ""
	if summary != "" {
		sumLine = fmt.Sprintf("\nKey points the video covers: %s", summary)
	}
	user := fmt.Sprintf(
		"Video title: %q%s%s\nLanguage: %s\n%s\n" +
			"Write the description in this order: 1) 1-2 sentence hook (with primary keyword), "+
			"2) 2-3 short paragraphs of context, 3) chapters block (if any), "+
			"4) hashtags (5-8, one line, no numbering), 5) CTA. "+
			"Total length: 150-220 words.",
		title, sumLine, tsLine, lang, ctaLine)
	return sys, user
}

// Thumbnail returns (system, user) prompts for the thumbnail concept generator.
func Thumbnail(title, mood string, faces int, colors, lang string) (string, string) {
	sys := "You are a YouTube thumbnail director. You design thumbnails with " +
		"high CTR by combining text overlays, expression, and visual contrast. " +
		"Reply in the user's language. Output the concept as a labeled list (one item per line): " +
		"'text overlay: ...', 'expression: ...', 'composition: ...', 'background: ...', " +
		"'image-gen prompt: ...' (a copy-pasteable prompt for Midjourney/DALL-E/Flux/SD)."
	colLine := ""
	if colors != "" {
		colLine = fmt.Sprintf("\nColor palette: %s", colors)
	}
	user := fmt.Sprintf(
		"Video title: %q\nMood: %s\nNumber of human faces: %d%s\nLanguage: %s\nGenerate a single, " +
			"specific thumbnail concept. Keep text overlay to 3-5 words max. " +
			"Make the image-gen prompt detailed (style, lighting, lens, composition).",
		title, mood, faces, colLine, lang)
	return sys, user
}

// Monetize returns (system, user) prompts for the monetization advisor.
func Monetize(niche string, subs int, region, focus, lang string) (string, string) {
	sys := "You are a YouTube monetization strategist. You give realistic, " +
		"tiered plans that match the channel's current size. Reply in the user's language. " +
		"Output as a labeled list (one item per line): 'tier: ...', 'estimated monthly: ...', " +
		"'first 3 steps: ...', 'pitfalls: ...'. No markdown fences, no preamble."
	user := fmt.Sprintf(
		"Niche: %q\nSubscribers (approx): %d\nPrimary region: %s\nFocus: %s\nLanguage: %s\n" +
			"Recommend a realistic monetization plan with 3 tiers (now / 6 months / 12 months). "+
			"Cover: ads RPM, sponsorships, digital products, affiliate, memberships/Patreon, "+
			"and any region-specific tips. Adjust aggressiveness to the subscriber count.",
		niche, subs, region, focus, lang)
	return sys, user
}

// Calendar returns (system, user) prompts for the content calendar planner.
func Calendar(month, niche, frequency, lang string) (string, string) {
	sys := "You are a YouTube content strategist. You design content calendars " +
		"that alternate formats and keep the audience engaged across the month. " +
		"Reply in the user's language. Output one planned video per line, format: " +
		"'Week X Day Y - Working title - Format (e.g. tutorial, vlog, shorts) - Hook angle'. " +
		"No markdown fences, no preamble."
	user := fmt.Sprintf(
		"Month: %s\nNiche: %q\nPosting frequency: %s\nLanguage: %s\n" +
			"Generate a full month content calendar. Mix formats (tutorial, shorts, vlog, "+
			"listicle, story). Make sure no two adjacent videos cover the same hook. "+
			"Order them by recommended publish day.",
		month, niche, frequency, lang)
	return sys, user
}
