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
	caasBaseURL       = "https://cataas.com/cat"
	caasSaysEndpoint  = "says"
	caasQueryStart    = "?"
	caasQueryAnd      = "&"
	caasReturnJSON    = "json=true"
	caasReturnHTML    = "html=true"
	caasPathSeparator = '/'
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

	//caasKeyJSON = "json"
	//caasKeyHTML = "html"

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
	ErrIDAndTag    = fmt.Errorf("cannot generate url with id and tag")
	ErrSaysNoText  = fmt.Errorf("cannot generate a Says URL with no text")
	ErrInvalidTag  = fmt.Errorf("invalid tag")
	ErrHTMLAndJSON = fmt.Errorf("cannot generate as both HTML and JSON")
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
	hasSays      bool // used to determine if using text overlay
	saysText     string
	customFilter bool
	params       []string // store params
	asJSON       bool
	asHTML       bool
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

func NewCatURL() *CatURL {
	return &CatURL{
		baseURL: caasBaseURL,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) WithTag(tag string) *CatURL {

	if !slices.Contains(AvailableTags, tag) {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}

	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          tag,
		hasTag:       true,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) WithSays(txt string) *CatURL {
	cleaned := url.PathEscape(txt)
	return &CatURL{
		baseURL:      c.baseURL,
		hasID:        c.hasID,
		catID:        c.catID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      true,
		saysText:     cleaned,
		customFilter: c.customFilter,
		params:       c.params,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: isCustom,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyFit, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyPosition, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyRed, strconv.Itoa(r))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyGreen, strconv.Itoa(g))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyBlue, strconv.Itoa(b))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}

	if !c.customFilter {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
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
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyBrightness, strconv.Itoa(brightness))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeySaturation, strconv.Itoa(saturation))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyHue, strconv.Itoa(hue))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyLightness, strconv.Itoa(lightness))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) WithFont(font CAASFont) *CatURL {
	if !c.hasSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyFont, str)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) WithFontSize(size int) *CatURL {
	if !c.hasSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyFontSize, strconv.Itoa(size))
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		hasID:        c.hasID,
		tag:          c.tag,
		hasTag:       c.hasTag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) WithFontColor(hexColor string) *CatURL {
	if !c.hasSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			hasID:        c.hasID,
			tag:          c.tag,
			hasTag:       c.hasTag,
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
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
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}

	updatedParams := c.updateParams(caasKeyFontColor, hexColor)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		tag:          c.tag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) WithFontBackground(hexColor string) *CatURL {
	if !c.hasSays {
		return &CatURL{
			baseURL:      c.baseURL,
			catID:        c.catID,
			tag:          c.tag,
			hasSays:      c.hasSays,
			saysText:     c.saysText,
			customFilter: c.customFilter,
			params:       c.params,
			asJSON:       c.asJSON,
			asHTML:       c.asHTML,
		}
	}
	updatedParams := c.updateParams(caasKeyFontBackground, hexColor)
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		tag:          c.tag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       updatedParams,
		asJSON:       c.asJSON,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) AsJSON() *CatURL {
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		tag:          c.tag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
		asJSON:       true,
		asHTML:       c.asHTML,
	}
}

func (c *CatURL) AsHTML() *CatURL {
	return &CatURL{
		baseURL:      c.baseURL,
		catID:        c.catID,
		tag:          c.tag,
		hasSays:      c.hasSays,
		saysText:     c.saysText,
		customFilter: c.customFilter,
		params:       c.params,
		asJSON:       c.asJSON,
		asHTML:       true,
	}
}

func (c *CatURL) Generate() (string, error) {

	// Bad Combos fail fast
	if c.hasID && c.hasTag {
		return "", ErrIDAndTag
	}
	if c.hasSays && c.saysText == "" {
		return "", ErrSaysNoText
	}
	if c.hasTag && !slices.Contains(AvailableTags, c.tag) {
		return "", ErrInvalidTag
	}
	if c.asHTML && c.asJSON {
		return "", ErrHTMLAndJSON
	}

	// write the base
	var b strings.Builder
	b.WriteString(c.baseURL)

	// add the ID/Tag if present
	if c.hasID {
		b.WriteRune(caasPathSeparator)
		b.WriteString(c.catID)
	}
	if c.hasTag {
		b.WriteRune(caasPathSeparator)
		b.WriteString(c.tag)
	}
	// add text overlay if present
	if c.hasSays {
		b.WriteRune(caasPathSeparator)
		b.WriteString(caasSaysEndpoint)
		b.WriteRune(caasPathSeparator)
		b.WriteString(c.saysText)
	}

	hasParams := len(c.params) > 0

	// if there are query params, process them
	if hasParams {
		b.WriteString(caasQueryStart)
		query := strings.Join(c.params, caasQueryAnd)
		b.WriteString(query)
	} else if c.asJSON || c.asHTML {
		b.WriteString(caasQueryStart)
	}
	// add output param if present
	if c.asJSON {
		b.WriteString(caasQueryAnd)
		b.WriteString(caasReturnJSON)
	}
	if c.asHTML {
		b.WriteString(caasQueryAnd)
		b.WriteString(caasReturnHTML)
	}

	return b.String(), nil

}

func ParseCatURL(catURL string) (*CatURL, error) {

	// let's validate it really is a URL
	parsedURL, err := url.Parse(catURL)
	if err != nil {
		return nil, err
	}

	// base url validation
	if parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("invalid URL scheme %s", parsedURL.Scheme)
	}
	if parsedURL.Host != "cataas.com" {
		return nil, fmt.Errorf("invalid URL host %s", parsedURL.Host)
	}
	if !strings.HasPrefix(parsedURL.Path, "/cat") {
		return nil, fmt.Errorf("invalid URL path %s", parsedURL.Path)
	}

	// we can get started now
	cu := NewCatURL()

	// split the path and validate the # parts
	pathParts := strings.Split(parsedURL.Path, "/")
	if len(pathParts) < 2 || len(pathParts) > 5 {
		return nil, fmt.Errorf("invalid URL path %s", parsedURL.Path)
	}

	// possible values for each index
	// 0 & 1 = "" & cat
	// 2 = says || id || tag
	// 3 = says || says-text
	// 4 = says-text

	// When we have 4 or 5, we must have a says where the last is text less than 4 cannot be says
	if len(pathParts) == 5 || len(pathParts) == 4 {
		cu.hasSays = true
		cu.saysText = pathParts[len(pathParts)-1]
	}

	//if we have 3 or 5 idx 2 is a tag or an id - if it's in available tags, it's a tag, otherwise must be an ID
	if (len(pathParts) == 3 || len(pathParts) == 5) && slices.Contains(AvailableTags, pathParts[2]) {
		cu.hasTag = true
		cu.tag = pathParts[2]
	} else if len(pathParts) == 3 || len(pathParts) == 5 {
		cu.hasID = true
		cu.catID = pathParts[2]
	}

	// Now the query
	q := parsedURL.Query()

	// split on &
	qp := strings.Split(parsedURL.RawQuery, "&")

	// check if custom filter applied
	if q.Get("filter") == "custom" {
		cu.customFilter = true
	}

	// check if the format is specified - it will be the last param, drop it
	if q.Has("json") {
		cu.asJSON = true
		qp = qp[:len(qp)-1]
	}
	if q.Has("html") {
		cu.asHTML = true
		qp = qp[:len(qp)-1]
	}

	// now we can tag our params
	cu.params = qp

	// and return
	return cu, nil
}
