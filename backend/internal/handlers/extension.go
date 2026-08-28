package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/yourorg/vscodemarket/internal/upstream"
)

type searchResp struct {
	Results []upstream.Extension `json:"results"`
	Total   int                  `json:"total"`
}

type versionsResp struct {
	Extension upstream.Extension         `json:"extension"`
	Versions  []upstream.ExtensionVersion `json:"versions"`
}

type versionResp struct {
	Extension upstream.Extension        `json:"extension"`
	Version   upstream.ExtensionVersion `json:"version"`
}

// ExtensionSearch handles GET /api/v1/extensions/search?q=&page=&pageSize=
func ExtensionSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	if q == "" {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "缺少查询参数 q", "missing query q"))
		return
	}
	page := atoiDefault(r.URL.Query().Get("page"), 1)
	pageSize := atoiDefault(r.URL.Query().Get("pageSize"), 10)
	results, total, err := upstream.SearchExtensions(r.Context(), q, page, pageSize)
	if err != nil {
		WriteError(w, http.StatusBadGateway, Err(CodeUpstreamFailure, "扩展搜索上游失败", "extension search upstream failure"))
		return
	}
	if results == nil {
		results = []upstream.Extension{}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(searchResp{Results: results, Total: total})
}

// ExtensionVersions handles GET /api/v1/extensions/{pub}/{name}/versions
func ExtensionVersions(w http.ResponseWriter, r *http.Request) {
	pub := strings.ToLower(strings.TrimSpace(r.PathValue("pub")))
	name := strings.TrimSpace(r.PathValue("name"))
	if pub == "" || name == "" {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "publisher 与 name 必填", "publisher and name required"))
		return
	}
	versions, err := upstream.FetchExtensionVersions(r.Context(), pub, name)
	if err != nil {
		if err == upstream.ErrExtensionNotFound {
			WriteError(w, http.StatusNotFound, Err(CodeExtensionVersionNF, "扩展不存在", "Extension not found"))
			return
		}
		WriteError(w, http.StatusBadGateway, Err(CodeUpstreamFailure, "扩展版本上游失败", "extension versions upstream failure"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(versionsResp{
		Extension: upstream.Extension{Publisher: pub, Name: name},
		Versions:  versions,
	})
}

// ExtensionVersion handles GET /api/v1/extensions/{pub}/{name}/versions/{ver}
func ExtensionVersion(w http.ResponseWriter, r *http.Request) {
	pub := strings.ToLower(strings.TrimSpace(r.PathValue("pub")))
	name := strings.TrimSpace(r.PathValue("name"))
	ver := strings.TrimSpace(r.PathValue("ver"))
	if pub == "" || name == "" || ver == "" {
		WriteError(w, http.StatusBadRequest, Err(CodeInvalidRequest, "publisher/name/ver 必填", "publisher/name/ver required"))
		return
	}
	versions, err := upstream.FetchExtensionVersions(r.Context(), pub, name)
	if err != nil {
		if err == upstream.ErrExtensionNotFound {
			WriteError(w, http.StatusNotFound, Err(CodeExtensionVersionNF, "扩展不存在", "Extension not found"))
			return
		}
		WriteError(w, http.StatusBadGateway, Err(CodeUpstreamFailure, "扩展版本上游失败", "extension versions upstream failure"))
		return
	}
	for _, v := range versions {
		if v.Version == ver {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_ = json.NewEncoder(w).Encode(versionResp{
				Extension: upstream.Extension{Publisher: pub, Name: name},
				Version:   v,
			})
			return
		}
	}
	WriteError(w, http.StatusNotFound, Err(CodeExtensionVersionNF, "扩展版本不存在", "Extension version not found"))
}

func atoiDefault(s string, d int) int {
	if s == "" {
		return d
	}
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return d
		}
		n = n*10 + int(c-'0')
	}
	return n
}