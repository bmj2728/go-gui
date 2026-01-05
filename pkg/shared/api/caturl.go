package api

import (
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/g4s8/hexcolor"
)

const (
	caasBaseURL      = "https://cataas.com/cat"
	caasSaysEndpoint = "says"
	caasQueryStart   = "?"
	caasQueryAnd     = "&"
	caasReturnJSON   = "json=true"
	caasReturnHTML   = "html=true"
)

const (
	// Basic Params

	caasKeyType     = "type"
	caasKeyFilter   = "filter"
	caasKeyFit      = "fit"
	caasKeyPosition = "position"
	caasKeyWidth    = "width"
	caasKeyHeight   = "height"
	caasKeyBlur     = "blur"

	// Custom Filter Params

	caasKeyRed        = "r"
	caasKeyGreen      = "g"
	caasKeyBlue       = "b"
	caasKeyBrightness = "brightness"
	caasKeySaturation = "saturation"
	caasKeyHue        = "hue"
	caasKeyLightness  = "lightness"

	// Font Params

	caasKeyFont           = "font"
	caasKeyFontSize       = "fontSize"
	caasKeyFontColor      = "fontColor"
	caasKeyFontBackground = "fontBackground"
)

var (
	ErrIDAndTag = fmt.Errorf("cannot generate url with id and tag")
)

func validRGBValue(val int) bool {
	if val < 0 || val > 255 {
		return false
	}
	return true
}

type CatURL struct {
	baseURL      string // will store the base url
	catID        string
	hasID        bool
	tag          string
	hasTag       bool
	isSays       bool // used to determine if using text overlay
	saysText     string
	customFilter bool
	params       []string // store params
}

func NewCatURL(baseURL string) *CatURL {
	return &CatURL{
		baseURL: baseURL,
		params:  make([]string, 0),
	}
}

func (c *CatURL) updateParams(key, value string) []string {
	param := fmt.Sprintf("%s=%s", key, value)
	updatedParams := append(c.params, param)
	return updatedParams
}

func (c *CatURL) WithID(id string) *CatURL {
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        id,
		hasID:        true,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
	}
}

func (c *CatURL) WithTag(tag string) *CatURL {

	if !slices.Contains(*AvailableTags, tag) {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}

	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          tag,
		hasTag:       true,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
	}
}

func (c *CatURL) WithSays(txt string) *CatURL {
	cleaned := url.QueryEscape(txt)
	return &CatURL{
		baseURL:      c.baseURL,
		hasID:        c.hasID,
		catID:        c.catID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       true,
		saysText:     cleaned,
		customFilter: c.customFilter,
		params:       c.params,
	}
}

func (c *CatURL) WithCAASImageType(imgType CAASImageType) *CatURL {
	// Get the str repr if it exists
	str, exists := CAASImageTypes[imgType]
	if !exists {
		// we're just returning the existing data
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}

	// generate the param
	updatedParams := c.updateParams(caasKeyType, str)

	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithCAASImageFilter(filter CAASImageFilter) *CatURL {
	str, exists := CAASImageFilters[filter]
	if !exists {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	isCustom := false
	if filter == CAASImageFilterCustom {
		isCustom = true
	}
	updatedParams := c.updateParams(caasKeyFilter, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: isCustom,
		params:       updatedParams,
	}
}

func (c *CatURL) WithCAASImageFit(fit CAASImageFit) *CatURL {
	str, exists := CAASImageFits[fit]
	if !exists {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyFit, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithCAASImagePosition(position CAASImagePosition) *CatURL {
	str, exists := CAASImagePositions[position]
	if !exists {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyPosition, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithWidth(width int) *CatURL {
	updatedParams := c.updateParams(caasKeyWidth, strconv.Itoa(width))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithHeight(height int) *CatURL {
	updatedParams := c.updateParams(caasKeyHeight, strconv.Itoa(height))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithBlur(blur int) *CatURL {
	updatedParams := c.updateParams(caasKeyBlur, strconv.Itoa(blur))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithFilterR(r int) *CatURL {
	if !validRGBValue(r) {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyRed, strconv.Itoa(r))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithFilterG(g int) *CatURL {
	if !validRGBValue(g) {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyGreen, strconv.Itoa(g))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithFilterB(b int) *CatURL {
	if !validRGBValue(b) {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyBlue, strconv.Itoa(b))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

// WithFilterRGB is a convenience function combining all 3 values
func (c *CatURL) WithFilterRGB(r, g, b int) *CatURL {
	if !validRGBValue(r) || !validRGBValue(g) || !validRGBValue(b) {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}

	if !c.customFilter {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}

	rParam := fmt.Sprintf("%s=%s", caasKeyRed, strconv.Itoa(r))
	gParam := fmt.Sprintf("%s=%s", caasKeyGreen, strconv.Itoa(g))
	bParam := fmt.Sprintf("%s=%s", caasKeyBlue, strconv.Itoa(b))

	updatedParams := append(c.params, rParam, gParam, bParam)

	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithBrightness(brightness int) *CatURL {
	if !c.customFilter {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyBrightness, strconv.Itoa(brightness))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithSaturation(saturation int) *CatURL {
	if !c.customFilter {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeySaturation, strconv.Itoa(saturation))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithHue(hue int) *CatURL {
	if !c.customFilter {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyHue, strconv.Itoa(hue))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithLightness(lightness int) *CatURL {
	if !c.customFilter {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyLightness, strconv.Itoa(lightness))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithFont(font CAASFont) *CatURL {
	if !c.isSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	str, exists := CAASFonts[font]
	if !exists {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	updatedParams := c.updateParams(caasKeyFont, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
	}
}

func (c *CatURL) WithFontSize(size int) *CatURL {
	if !c.isSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	c.updateParams(caasKeyFontSize, strconv.Itoa(size))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
	}
}

func (c *CatURL) WithFontColor(hexColor string) *CatURL {
	if !c.isSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}

	_, err := hexcolor.Parse(hexColor)
	if err != nil {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}

	c.updateParams(caasKeyFontColor, hexColor)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		tag:          c.tag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
	}
}

func (c *CatURL) WithFontBackground(hexColor string) *CatURL {
	if !c.isSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			tag:          c.tag,
			isSays:       c.isSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
		}
	}
	c.updateParams(caasKeyFontBackground, hexColor)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		tag:          c.tag,
		isSays:       c.isSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
	}
}

func (c *CatURL) Generate() (string, error) {

	// Bad Combo fail fast
	if c.hasID && c.hasTag {
		return "", ErrIDAndTag
	}

	// write the base
	var b strings.Builder
	b.WriteString(c.baseURL)

	// add the ID/Tag if present
	if c.hasID {
		b.WriteRune('/')
		b.WriteString(c.catID)
	}
	if c.hasTag {
		b.WriteRune('/')
		b.WriteString(c.tag)
	}

	if c.isSays {
		b.WriteRune('/')
		b.WriteString(caasSaysEndpoint)
		b.WriteRune('/')
		b.WriteString(c.saysText)
	}

	if len(c.params) > 0 {
		b.WriteString(caasQueryStart)
		query := strings.Join(c.params, caasQueryAnd)
		b.WriteString(query)
	}

	return b.String(), nil

}

/*
	valid endpoints:
	basicCalls:
		- caasBaseURL
		- caasBaseURL/ID
		- caasBaseURL/TAG
	withTextOverlay:
		- caasBaseURL/caasSaysEndpoint/escaped%20text%21
		- caasBaseURL/ID/caasSaysEndpoint/escaped%20text%21
		- caasBaseURL/TAG/caasSaysEndpoint/escaped%20text%21
	queryCalls:
		- caasBaseURL + caasQueryStart + params
		- caasBaseURL/ID + caasQueryStart + params
		- caasBaseURL/TAG + caasQueryStart + params
		- caasBaseURL/caasSaysEndpoint/escaped%20text%21 + caasQueryStart + params
		- caasBaseURL/ID/caasSaysEndpoint/escaped%20text%21 + caasQueryStart + params
		- caasBaseURL/TAG/caasSaysEndpoint/escaped%20text%21 + caasQueryStart + params
*/
