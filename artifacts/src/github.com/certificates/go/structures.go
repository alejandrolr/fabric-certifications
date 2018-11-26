/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

// constants
const ISSUER_LIST = "issuer-list"
const BADGE_PREFIX = "badge:"
const CERT_PREFIX = "cert:"

// Issuer storage structures
type IssuerList struct {
	IssuerSummary []IssuerSummary `json:"issuers"`
}

type IssuerSummary struct {
	Email    string   `json:"issuerEmail"`
	BagdeIDs []string `json:"badgeIDs"`
	CertIDs  []string `json:"certIDs"`
}

// Receiver storage structures
type ReceiverList struct {
	ReceiverSummary []ReceiverSummary `json:"receivers"`
}

type ReceiverSummary struct {
	Email    string   `json:"receiverEmail"`
	CertsIDs []string `json:"certsIDs"`
}

// Certificate structures
type Certificate struct {
	Context          string           `json:"@context"`
	Id               string           `json:"id"`
	Type             string           `json:"type"`
	IssuedOn         string           `json:"issuedOn"`
	Recipient        Recipient        `json:"recipient"`
	RecipientProfile RecipientProfile `json:"recipientProfile"`
	Verification     Verification     `json:"verification"`
	Badge            Badge            `json:"badge"`
}

type Recipient struct {
	Identity string `json:"identity"`
	Type     string `json:"type"`
	Hashed   bool   `json:"hashed"`
}

type RecipientProfile struct {
	PublicKey string   `json:"publicKey"`
	Name      string   `json:"name"`
	Type      []string `json:"type"`
}

type Verification struct {
	Location string   `json:"location"`
	Type     []string `json:"type"`
}

type Badge struct {
	Id             string           `json:"id"`
	Name           string           `json:"name"`
	Type           string           `json:"type"`
	Issuer         Issuer           `json:"issuer"`
	Criteria       Criteria         `json:"criteria"`
	Image          string           `json:"image"`
	Description    string           `json:"description"`
	SignatureLines []SignatureLines `json:"signatureLines"`
}

type Issuer struct {
	Id    string `json:"id"`
	Url   string `json:"url"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Type  string `json:"type"`
	// Image          string `json:"image"`
	RevocationList string `json:"revocationList"`
}

type Criteria struct {
	Narrative string `json:"narrative"`
}

type SignatureLines struct {
	JobTitle string   `json:"jobTitle"`
	Name     string   `json:"name"`
	Type     []string `json:"type"`
	Image    string   `json:"image"`
}
