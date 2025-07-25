// SPDX-FileCopyrightText: Copyright The Miniflux Authors. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package urllib // import "miniflux.app/v2/internal/urllib"

import (
	"fmt"
	"net/url"
	"strings"
)

// IsAbsoluteURL returns true if the link is absolute.
func IsAbsoluteURL(link string) bool {
	u, err := url.Parse(link)
	if err != nil {
		return false
	}
	return u.IsAbs()
}

// GetAbsoluteURL return the absolute form of `input` is possible, as well as its parser form.
func GetAbsoluteURL(input string) (string, *url.URL, error) {
	if strings.HasPrefix(input, "//") {
		return "https:" + input, nil, nil
	}
	if strings.HasPrefix(input, "https://") || strings.HasPrefix(input, "http://") {
		return input, nil, nil
	}

	u, err := url.Parse(input)
	if err != nil {
		return "", nil, fmt.Errorf("unable to parse input URL: %v", err)
	}

	if u.IsAbs() {
		return u.String(), u, nil
	}
	return "", u, nil
}

// AbsoluteURL converts the input URL as absolute URL if necessary.
func AbsoluteURL(baseURL, input string) (string, error) {
	absURL, u, err := GetAbsoluteURL(input)
	if err != nil {
		return "", err
	}
	if absURL != "" {
		return absURL, nil
	}

	base, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("unable to parse base URL: %v", err)
	}

	return base.ResolveReference(u).String(), nil
}

// RootURL returns absolute URL without the path.
func RootURL(websiteURL string) string {
	if strings.HasPrefix(websiteURL, "//") {
		websiteURL = "https://" + websiteURL[2:]
	}

	absoluteURL, err := AbsoluteURL(websiteURL, "")
	if err != nil {
		return websiteURL
	}

	u, err := url.Parse(absoluteURL)
	if err != nil {
		return absoluteURL
	}

	return u.Scheme + "://" + u.Host + "/"
}

// IsHTTPS returns true if the URL is using HTTPS.
func IsHTTPS(websiteURL string) bool {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return false
	}

	return strings.EqualFold(parsedURL.Scheme, "https")
}

// Domain returns only the domain part of the given URL.
func Domain(websiteURL string) string {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return websiteURL
	}

	return parsedURL.Host
}

// DomainWithoutWWW returns only the domain part of the given URL, with the "www." prefix removed if present.
func DomainWithoutWWW(websiteURL string) string {
	return strings.TrimPrefix(Domain(websiteURL), "www.")
}

// JoinBaseURLAndPath returns a URL string with the provided path elements joined together.
func JoinBaseURLAndPath(baseURL, path string) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("empty base URL")
	}

	if path == "" {
		return "", fmt.Errorf("empty path")
	}

	_, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	finalURL, err := url.JoinPath(baseURL, path)
	if err != nil {
		return "", fmt.Errorf("unable to join base URL %s and path %s: %w", baseURL, path, err)
	}

	return finalURL, nil
}
