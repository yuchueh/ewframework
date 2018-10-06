package context

import (
	"path"
	"regexp"
	"strings"
	"github.com/yuchueh/ewframework/utils"
)

var (
	allowSuffixExt = []string{".json", ".xml", ".html"}
)

//Wildcard 通配符
// Tree has three elements: FixRouter/Wildcard/Leaves
// fixRouter sotres Fixed Router
// Wildcard stores params
// Leaves store the endpoint information
type Tree struct {
	//Prefix set for static router
	Prefix string
	//search fix route first
	Fixrouters []*Tree
	//if set, failure to match Fixrouters search then search Wildcard
	Wildcard *Tree
	//if set, failure to match Wildcard search
	Leaves []*leafInfo
}

// NewTree return a new Tree
func NewTree() *Tree {
	return &Tree{}
}

// AddTree will add tree to the exist Tree
// Prefix should has no params
func (t *Tree) AddTree(Prefix string, tree *Tree) {
	t.addtree(splitPath(Prefix), tree, nil, "")
}

func (t *Tree) addtree(segments []string, tree *Tree, Wildcards []string, reg string) {
	if len(segments) == 0 {
		panic("Prefix should has path")
	}
	seg := segments[0]
	iswild, params, regexpStr := splitSegment(seg)
	// if it's ? meaning can igone this, so add one more rule for it
	if len(params) > 0 && params[0] == ":" {
		params = params[1:]
		if len(segments[1:]) > 0 {
			t.addtree(segments[1:], tree, append(Wildcards, params...), reg)
		} else {
			filterTreeWithPrefix(tree, Wildcards, reg)
		}
	}
	//Rule: /login/*/access match /login/2009/11/access
	//if already has *, and when loop the access, should as a regexpStr
	if !iswild && utils.InSlice(":splat", Wildcards) {
		iswild = true
		regexpStr = seg
	}
	//Rule: /user/:id/*
	if seg == "*" && len(Wildcards) > 0 && reg == "" {
		regexpStr = "(.+)"
	}
	if len(segments) == 1 {
		if iswild {
			if regexpStr != "" {
				if reg == "" {
					rr := ""
					for _, w := range Wildcards {
						if w == ":splat" {
							rr = rr + "(.+)/"
						} else {
							rr = rr + "([^/]+)/"
						}
					}
					regexpStr = rr + regexpStr
				} else {
					regexpStr = "/" + regexpStr
				}
			} else if reg != "" {
				if seg == "*.*" {
					regexpStr = "([^.]+).(.+)"
				} else {
					for _, w := range params {
						if w == "." || w == ":" {
							continue
						}
						regexpStr = "([^/]+)/" + regexpStr
					}
				}
			}
			reg = strings.Trim(reg+"/"+regexpStr, "/")
			filterTreeWithPrefix(tree, append(Wildcards, params...), reg)
			t.Wildcard = tree
		} else {
			reg = strings.Trim(reg+"/"+regexpStr, "/")
			filterTreeWithPrefix(tree, append(Wildcards, params...), reg)
			tree.Prefix = seg
			t.Fixrouters = append(t.Fixrouters, tree)
		}
		return
	}

	if iswild {
		if t.Wildcard == nil {
			t.Wildcard = NewTree()
		}
		if regexpStr != "" {
			if reg == "" {
				rr := ""
				for _, w := range Wildcards {
					if w == ":splat" {
						rr = rr + "(.+)/"
					} else {
						rr = rr + "([^/]+)/"
					}
				}
				regexpStr = rr + regexpStr
			} else {
				regexpStr = "/" + regexpStr
			}
		} else if reg != "" {
			if seg == "*.*" {
				regexpStr = "([^.]+).(.+)"
				params = params[1:]
			} else {
				for range params {
					regexpStr = "([^/]+)/" + regexpStr
				}
			}
		} else {
			if seg == "*.*" {
				params = params[1:]
			}
		}
		reg = strings.TrimRight(strings.TrimRight(reg, "/")+"/"+regexpStr, "/")
		t.Wildcard.addtree(segments[1:], tree, append(Wildcards, params...), reg)
	} else {
		subTree := NewTree()
		subTree.Prefix = seg
		t.Fixrouters = append(t.Fixrouters, subTree)
		subTree.addtree(segments[1:], tree, append(Wildcards, params...), reg)
	}
}

func filterTreeWithPrefix(t *Tree, Wildcards []string, reg string) {
	for _, v := range t.Fixrouters {
		filterTreeWithPrefix(v, Wildcards, reg)
	}
	if t.Wildcard != nil {
		filterTreeWithPrefix(t.Wildcard, Wildcards, reg)
	}
	for _, l := range t.Leaves {
		if reg != "" {
			if l.Regexps != nil {
				l.Wildcards = append(Wildcards, l.Wildcards...)
				l.Regexps = regexp.MustCompile("^" + reg + "/" + strings.Trim(l.Regexps.String(), "^$") + "$")
			} else {
				for _, v := range l.Wildcards {
					if v == ":splat" {
						reg = reg + "/(.+)"
					} else {
						reg = reg + "/([^/]+)"
					}
				}
				l.Regexps = regexp.MustCompile("^" + reg + "$")
				l.Wildcards = append(Wildcards, l.Wildcards...)
			}
		} else {
			l.Wildcards = append(Wildcards, l.Wildcards...)
			if l.Regexps != nil {
				for _, w := range Wildcards {
					if w == ":splat" {
						reg = "(.+)/" + reg
					} else {
						reg = "([^/]+)/" + reg
					}
				}
				l.Regexps = regexp.MustCompile("^" + reg + strings.Trim(l.Regexps.String(), "^$") + "$")
			}
		}
	}
}

// AddRouter call addseg function
func (t *Tree) AddRouter(pattern string, RunObject interface{}) {
	t.addseg(splitPath(pattern), RunObject, nil, "")
}

// "/"
// "admin" ->
func (t *Tree) addseg(segments []string, route interface{}, Wildcards []string, reg string) {
	if len(segments) == 0 {
		if reg != "" {
			t.Leaves = append(t.Leaves, &leafInfo{RunObject: route, Wildcards: Wildcards, Regexps: regexp.MustCompile("^" + reg + "$")})
		} else {
			t.Leaves = append(t.Leaves, &leafInfo{RunObject: route, Wildcards: Wildcards})
		}
	} else {
		seg := segments[0]
		iswild, params, regexpStr := splitSegment(seg)
		// if it's ? meaning can igone this, so add one more rule for it
		if len(params) > 0 && params[0] == ":" {
			t.addseg(segments[1:], route, Wildcards, reg)
			params = params[1:]
		}
		//Rule: /login/*/access match /login/2009/11/access
		//if already has *, and when loop the access, should as a regexpStr
		if !iswild && utils.InSlice(":splat", Wildcards) {
			iswild = true
			regexpStr = seg
		}
		//Rule: /user/:id/*
		if seg == "*" && len(Wildcards) > 0 && reg == "" {
			regexpStr = "(.+)"
		}
		if iswild {
			if t.Wildcard == nil {
				t.Wildcard = NewTree()
			}
			if regexpStr != "" {
				if reg == "" {
					rr := ""
					for _, w := range Wildcards {
						if w == ":splat" {
							rr = rr + "(.+)/"
						} else {
							rr = rr + "([^/]+)/"
						}
					}
					regexpStr = rr + regexpStr
				} else {
					regexpStr = "/" + regexpStr
				}
			} else if reg != "" {
				if seg == "*.*" {
					regexpStr = "/([^.]+).(.+)"
					params = params[1:]
				} else {
					for range params {
						regexpStr = "/([^/]+)" + regexpStr
					}
				}
			} else {
				if seg == "*.*" {
					params = params[1:]
				}
			}
			t.Wildcard.addseg(segments[1:], route, append(Wildcards, params...), reg+regexpStr)
		} else {
			var subTree *Tree
			for _, sub := range t.Fixrouters {
				if sub.Prefix == seg {
					subTree = sub
					break
				}
			}
			if subTree == nil {
				subTree = NewTree()
				subTree.Prefix = seg
				t.Fixrouters = append(t.Fixrouters, subTree)
			}
			subTree.addseg(segments[1:], route, Wildcards, reg)
		}
	}
}

// Match router to RunObject & params
func (t *Tree) Match(pattern string, ctx *Context) (RunObject interface{}) {
	if len(pattern) == 0 || pattern[0] != '/' {
		return nil
	}
	w := make([]string, 0, 20)
	return t.match(pattern[1:], pattern, w, ctx)
}

func (t *Tree) match(treePattern string, pattern string, WildcardValues []string, ctx *Context) (RunObject interface{}) {
	if len(pattern) > 0 {
		i := 0
		for ; i < len(pattern) && pattern[i] == '/'; i++ {
		}
		pattern = pattern[i:]
	}
	// Handle leaf nodes:
	if len(pattern) == 0 {
		for _, l := range t.Leaves {
			if ok := l.match(treePattern, WildcardValues, ctx); ok {
				return l.RunObject
			}
		}
		if t.Wildcard != nil {
			for _, l := range t.Wildcard.Leaves {
				if ok := l.match(treePattern, WildcardValues, ctx); ok {
					return l.RunObject
				}
			}
		}
		return nil
	}
	var seg string
	i, l := 0, len(pattern)
	for ; i < l && pattern[i] != '/'; i++ {
	}
	if i == 0 {
		seg = pattern
		pattern = ""
	} else {
		seg = pattern[:i]
		pattern = pattern[i:]
	}
	for _, subTree := range t.Fixrouters {
		if subTree.Prefix == seg {
			if len(pattern) != 0 && pattern[0] == '/' {
				treePattern = pattern[1:]
			} else {
				treePattern = pattern
			}
			RunObject = subTree.match(treePattern, pattern, WildcardValues, ctx)
			if RunObject != nil {
				break
			}
		}
	}
	if RunObject == nil && len(t.Fixrouters) > 0 {
		// Filter the .json .xml .html extension
		for _, str := range allowSuffixExt {
			if strings.HasSuffix(seg, str) {
				for _, subTree := range t.Fixrouters {
					if subTree.Prefix == seg[:len(seg)-len(str)] {
						RunObject = subTree.match(treePattern, pattern, WildcardValues, ctx)
						if RunObject != nil {
							ctx.Input.SetParam(":ext", str[1:])
						}
					}
				}
			}
		}
	}
	if RunObject == nil && t.Wildcard != nil {
		RunObject = t.Wildcard.match(treePattern, pattern, append(WildcardValues, seg), ctx)
	}

	if RunObject == nil && len(t.Leaves) > 0 {
		WildcardValues = append(WildcardValues, seg)
		start, i := 0, 0
		for ; i < len(pattern); i++ {
			if pattern[i] == '/' {
				if i != 0 && start < len(pattern) {
					WildcardValues = append(WildcardValues, pattern[start:i])
				}
				start = i + 1
				continue
			}
		}
		if start > 0 {
			WildcardValues = append(WildcardValues, pattern[start:i])
		}
		for _, l := range t.Leaves {
			if ok := l.match(treePattern, WildcardValues, ctx); ok {
				return l.RunObject
			}
		}
	}
	return RunObject
}

type leafInfo struct {
	// names of Wildcards that lead to this leaf. eg, ["id" "name"] for the Wildcard ":id" and ":name"
	Wildcards []string

	// if the leaf is regexp
	Regexps *regexp.Regexp

	RunObject interface{}
}

func (leaf *leafInfo) match(treePattern string, WildcardValues []string, ctx *Context) (ok bool) {
	//fmt.Println("Leaf:", WildcardValues, leaf.Wildcards, leaf.Regexps)
	if leaf.Regexps == nil {
		if len(WildcardValues) == 0 && len(leaf.Wildcards) == 0 { // static path
			return true
		}
		// match *
		if len(leaf.Wildcards) == 1 && leaf.Wildcards[0] == ":splat" {
			ctx.Input.SetParam(":splat", treePattern)
			return true
		}
		// match *.* or :id
		if len(leaf.Wildcards) >= 2 && leaf.Wildcards[len(leaf.Wildcards)-2] == ":path" && leaf.Wildcards[len(leaf.Wildcards)-1] == ":ext" {
			if len(leaf.Wildcards) == 2 {
				lastone := WildcardValues[len(WildcardValues)-1]
				strs := strings.SplitN(lastone, ".", 2)
				if len(strs) == 2 {
					ctx.Input.SetParam(":ext", strs[1])
				}
				ctx.Input.SetParam(":path", path.Join(path.Join(WildcardValues[:len(WildcardValues)-1]...), strs[0]))
				return true
			} else if len(WildcardValues) < 2 {
				return false
			}
			var index int
			for index = 0; index < len(leaf.Wildcards)-2; index++ {
				ctx.Input.SetParam(leaf.Wildcards[index], WildcardValues[index])
			}
			lastone := WildcardValues[len(WildcardValues)-1]
			strs := strings.SplitN(lastone, ".", 2)
			if len(strs) == 2 {
				ctx.Input.SetParam(":ext", strs[1])
			}
			if index > (len(WildcardValues) - 1) {
				ctx.Input.SetParam(":path", "")
			} else {
				ctx.Input.SetParam(":path", path.Join(path.Join(WildcardValues[index:len(WildcardValues)-1]...), strs[0]))
			}
			return true
		}
		// match :id
		if len(leaf.Wildcards) != len(WildcardValues) {
			return false
		}
		for j, v := range leaf.Wildcards {
			ctx.Input.SetParam(v, WildcardValues[j])
		}
		return true
	}

	if !leaf.Regexps.MatchString(path.Join(WildcardValues...)) {
		return false
	}
	matches := leaf.Regexps.FindStringSubmatch(path.Join(WildcardValues...))
	for i, match := range matches[1:] {
		if i < len(leaf.Wildcards) {
			ctx.Input.SetParam(leaf.Wildcards[i], match)
		}
	}
	return true
}

// "/" -> []
// "/admin" -> ["admin"]
// "/admin/" -> ["admin"]
// "/admin/users" -> ["admin", "users"]
func splitPath(key string) []string {
	key = strings.Trim(key, "/ ")
	if key == "" {
		return []string{}
	}
	return strings.Split(key, "/")
}

// "admin" -> false, nil, ""
// ":id" -> true, [:id], ""
// "?:id" -> true, [: :id], ""        : meaning can empty
// ":id:int" -> true, [:id], ([0-9]+)
// ":name:string" -> true, [:name], ([\w]+)
// ":id([0-9]+)" -> true, [:id], ([0-9]+)
// ":id([0-9]+)_:name" -> true, [:id :name], ([0-9]+)_(.+)
// "cms_:id_:page.html" -> true, [:id_ :page], cms_(.+)(.+).html
// "cms_:id(.+)_:page.html" -> true, [:id :page], cms_(.+)_(.+).html
// "*" -> true, [:splat], ""
// "*.*" -> true,[. :path :ext], ""      . meaning separator
func splitSegment(key string) (bool, []string, string) {
	if strings.HasPrefix(key, "*") {
		if key == "*.*" {
			return true, []string{".", ":path", ":ext"}, ""
		}
		return true, []string{":splat"}, ""
	}
	if strings.ContainsAny(key, ":") {
		var paramsNum int
		var out []rune
		var start bool
		var startexp bool
		var param []rune
		var expt []rune
		var skipnum int
		params := []string{}
		reg := regexp.MustCompile(`[a-zA-Z0-9_]+`)
		for i, v := range key {
			if skipnum > 0 {
				skipnum--
				continue
			}
			if start {
				//:id:int and :name:string
				if v == ':' {
					if len(key) >= i+4 {
						if key[i+1:i+4] == "int" {
							out = append(out, []rune("([0-9]+)")...)
							params = append(params, ":"+string(param))
							start = false
							startexp = false
							skipnum = 3
							param = make([]rune, 0)
							paramsNum++
							continue
						}
					}
					if len(key) >= i+7 {
						if key[i+1:i+7] == "string" {
							out = append(out, []rune(`([\w]+)`)...)
							params = append(params, ":"+string(param))
							paramsNum++
							start = false
							startexp = false
							skipnum = 6
							param = make([]rune, 0)
							continue
						}
					}
				}
				// params only support a-zA-Z0-9
				if reg.MatchString(string(v)) {
					param = append(param, v)
					continue
				}
				if v != '(' {
					out = append(out, []rune(`(.+)`)...)
					params = append(params, ":"+string(param))
					param = make([]rune, 0)
					paramsNum++
					start = false
					startexp = false
				}
			}
			if startexp {
				if v != ')' {
					expt = append(expt, v)
					continue
				}
			}
			// Escape Sequence '\'
			if i > 0 && key[i-1] == '\\' {
				out = append(out, v)
			} else if v == ':' {
				param = make([]rune, 0)
				start = true
			} else if v == '(' {
				startexp = true
				start = false
				if len(param) > 0 {
					params = append(params, ":"+string(param))
					param = make([]rune, 0)
				}
				paramsNum++
				expt = make([]rune, 0)
				expt = append(expt, '(')
			} else if v == ')' {
				startexp = false
				expt = append(expt, ')')
				out = append(out, expt...)
				param = make([]rune, 0)
			} else if v == '?' {
				params = append(params, ":")
			} else {
				out = append(out, v)
			}
		}
		if len(param) > 0 {
			if paramsNum > 0 {
				out = append(out, []rune(`(.+)`)...)
			}
			params = append(params, ":"+string(param))
		}
		return true, params, string(out)
	}
	return false, nil, ""
}
