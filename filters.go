package gendoc

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	paraPattern         = regexp.MustCompile(`(\n|\r|\r\n)\s*`)
	spacePattern        = regexp.MustCompile("( )+")
	multiNewlinePattern = regexp.MustCompile(`(\r\n|\r|\n){2,}`)
	basename            = "v2ray.core."
)

// PFilter splits the content by new lines and wraps each one in a <p> tag.
func PFilter(content string) template.HTML {
	paragraphs := paraPattern.Split(content, -1)
	return template.HTML(fmt.Sprintf("<p>%s</p>", strings.Join(paragraphs, "</p><p>")))
}

// ParaFilter splits the content by new lines and wraps each one in a <para> tag.
func ParaFilter(content string) string {
	paragraphs := paraPattern.Split(content, -1)
	return fmt.Sprintf("<para>%s</para>", strings.Join(paragraphs, "</para><para>"))
}

// NoBrFilter removes single CR and LF from content.
func NoBrFilter(content string) string {
	normalized := strings.Replace(content, "\r\n", "\n", -1)
	paragraphs := multiNewlinePattern.Split(normalized, -1)
	for i, p := range paragraphs {
		withoutCR := strings.Replace(p, "\r", " ", -1)
		withoutLF := strings.Replace(withoutCR, "\n", " ", -1)
		paragraphs[i] = spacePattern.ReplaceAllString(withoutLF, " ")
	}
	return strings.Join(paragraphs, "\n\n")
}

func ImportLinkFilter(content []string) []string {
	m := make(map[string]bool, len(content))
	for _, e := range content {
		m[path.Dir(path.Clean("/"+os.Getenv("PROTOC_GEN_DOC_PAGE_ROOT")+"/"+e))] = true
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, strings.TrimPrefix(k, "/"))
	}
	return keys
}

func TypeLinkParser(parent string, this_type string, full_type string, long_type string, this_pack string) string {
	if !strings.Contains(full_type, basename) {
		return "#" + full_type
	}

	// println("p", parent)
	// println("this", this_type)
	// println("full", full_type)
	// println("long", long_type)
	// println("pack", this_pack)
	// println(">")

	isInCurrentPackage := this_pack+"."+this_type == full_type
	if isInCurrentPackage {
		return "#" + full_type
	}

	isNestedSubMessage := parent+"."+this_type == full_type
	if isNestedSubMessage {
		return "#" + full_type
	}

	isInSubPackage := this_pack+"."+long_type == full_type
	if isInSubPackage {
		tmp := strings.TrimSuffix(long_type, "."+this_type)
		long_type_path := strings.ReplaceAll(tmp, ".", "/")
		return long_type_path + "/index.html#" + full_type
	}

	tmp := strings.TrimPrefix(full_type, basename)
	tmp = strings.ReplaceAll(tmp, ".", "/")
	tmp = path.Dir(tmp)

	return "/" + tmp + "/index.html#" + full_type
}

func commonPackage(files []*File) string {
	for _, f := range files {
		return f.Package
	}
	return "UNKNOWN PACKAGE"
}
