package emailx

func InviteUserEmail(inviteLink string) (subject, htmlBody, textBody string) {
	return "Confirm your Zula account",
		"<h1>Start Reading with Zula!</h1><p>Hi there!</p><p>Click the link below to confirm your account.</p> " +
			"<a href='" + inviteLink + "'>Confirmation link</a><br><p>The Reading Tribe</p>",
		"Start Reading with Zula. Hi there! Click the link below to confirm your account. " + inviteLink
}
