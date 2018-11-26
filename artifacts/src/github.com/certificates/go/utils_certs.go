package main

func createIssuerSummary(issuer Issuer) IssuerSummary {
	issuerSummary := IssuerSummary{
		Email: issuer.Email,
		// Blank parameters by default
	}
	return issuerSummary
}

// Function that creates certificate
func createCertificate(id, issuedOn string, rec Recipient, recProf RecipientProfile,
	ver Verification, badge Badge) Certificate {

	logger.Infof("Creating certificate")
	certificate := Certificate{
		Context:          "https://w3id.org/openbadges/v2",
		IssuedOn:         issuedOn,
		Id:               id,
		Type:             "Assertion",
		Recipient:        rec,
		RecipientProfile: recProf,
		Verification:     ver,
		Badge:            badge,
	}

	return certificate
}

func createRecipient(identity string) Recipient {
	recipient := Recipient{
		Identity: identity,
		Type:     "email",
		Hashed:   false,
	}
	return recipient
}

func createRecipientProfile(pubKey, name string) RecipientProfile {
	recipientProfile := RecipientProfile{
		PublicKey: pubKey,
		Name:      name,
		Type:      []string{"RecipientProfile", "Extension"},
	}
	return recipientProfile
}

func createVerification(location string) Verification {
	verification := Verification{
		Location: location,
		Type:     []string{"MerkleProofVerification2017", "Extension"},
	}
	return verification
}

func createBadge(id, name, description string, issuer Issuer, criteria Criteria, signature SignatureLines) Badge {

	badge := Badge{
		Id:             id,
		Name:           name,
		Type:           "BadgeClass",
		Issuer:         createIssuer(issuer.Id, issuer.Url, issuer.Email, issuer.Name),
		Criteria:       createCriteria(criteria.Narrative),
		Description:    description,
		SignatureLines: []SignatureLines{createSignatureLines(signature.JobTitle, signature.Name)},
	}
	return badge
}

func createIssuer(id, url, email, name string) Issuer {

	issuer := Issuer{
		Id:             id,
		Url:            url,
		Name:           name,
		Email:          email,
		Type:           "Profile",
		RevocationList: "",
	}
	return issuer
}

func createCriteria(narrative string) Criteria {
	criteria := Criteria{
		Narrative: narrative,
	}
	return criteria
}

func createSignatureLines(jobTitle, name string) SignatureLines {
	signatureLines := SignatureLines{
		JobTitle: jobTitle,
		Name:     name,
		Type:     []string{"SignatureLine", "Extension"},
	}
	return signatureLines
}
