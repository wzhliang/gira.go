package gitee

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/wzhliang/gira/pkg/context"
)

const apiRoot = "https://gitee.com/api/v5/repos/%s/%s"

func mapToJSONString(data map[string]string) string {
	js, err := json.Marshal(data)
	if err != nil {
		fmt.Println(data)
		panic("invalid data")
	}
	return string(js)
}

// Gitee ...
type Gitee struct {
	User  string
	Token string
	root  string
}

// New ...
func New(repo, owner, token string) *Gitee {
	log.Info().Msgf("Gitee::New %s, %s, %s", repo, owner, token)
	return &Gitee{
		User:  owner,
		Token: token,
		root:  fmt.Sprintf(apiRoot, owner, repo),
	}
}

func findIssueID(title string) string {
	return strings.Split(title, " ")[0]
}

// GetPR ...
func (g *Gitee) GetPR(ctx *context.Context) error {
	log.Info().Msgf("Gitee::GetPR %s ...", ctx.PR.ID)
	addr := fmt.Sprintf("%s/pulls/%s", g.root, ctx.PR.ID)

	resp, err := g.get(addr)
	if err != nil {
		log.Error().Msg("Failed to GET from gitee server")
		return nil
	}

	var data PR
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response")
	}

	// FIXME: not sure if this is the best way to do it but seems clean enough
	ctx.PR.Title = data.Title
	ctx.PR.URL = data.URL
	ctx.Issue.ID = findIssueID(data.Title)
	for _, u := range data.Assignees {
		ctx.PR.Owners = append(ctx.PR.Owners, u.Name)
	}
	for _, u := range data.Testers {
		ctx.PR.Owners = append(ctx.PR.Owners, u.Name)
	}

	return nil
}

// CreatePR ...
func (g *Gitee) CreatePR(ctx *context.Context) error {
	log.Info().Msg("Creating PR...")
	addr := fmt.Sprintf("%s/pulls", g.root)
	resp, err := g.post(addr, map[string]string{
		"title": fmt.Sprintf("%s %s", ctx.Issue.ID, ctx.Issue.Summary),
		"head":  ctx.Issue.ID,
		"base":  ctx.PR.TargetBranch,
		"body":  fmt.Sprintf("PR for %s", ctx.Issue.URL),
	})

	var data PR
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response")
	}
	ctx.PR.URL = data.URL

	return err
}

// MergePR ...
func (g *Gitee) MergePR(ctx *context.Context) error {
	log.Printf("Merging PR...")
	addr := fmt.Sprintf("%s/pulls/%s/merge", g.root, ctx.PR.ID)
	resp, err := g.put(addr, map[string]string{})
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 300 {
		return errors.New("failed")
	}
	return nil
}

// GetBranch ...
func (g *Gitee) GetBranch(ctx *context.Context) error {
	log.Info().Msgf("Gitee::GetBranch %s", ctx.Issue.ID)
	addr := fmt.Sprintf("%s/branches/%s", g.root, ctx.Issue.ID) // issue == branch

	resp, err := g.get(addr)
	if err != nil {
		log.Error().Msg("Failed to GET from gitee server")
		return nil
	}

	if resp.StatusCode() == 200 {
		return nil
	} else if resp.StatusCode() == 404 {
		return errors.New("branch doesn't exits")
	}
	return errors.New("error getting branch")
}

// List ...
func (g *Gitee) List(ctx *context.Context) ([]string, error) {
	log.Info().Msg("ListPRs ...")

	resp, err := g.get(fmt.Sprintf("%s/pulls", g.root))
	if err != nil {
		return nil, nil
	}

	var data []PR
	err = json.Unmarshal(resp.Body(), &data)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response")
	}

	var ret []string
	for _, pr := range data {
		ret = append(ret, fmt.Sprintf("%d - %s", pr.Number, pr.Title))
	}

	return ret, nil
}

func (g *Gitee) get(addr string) (*resty.Response, error) {
	log.Info().Msgf("addr: %s", addr)
	client := resty.New()
	client.SetQueryParam("access_token", g.Token)
	resp, err := client.R().
		EnableTrace().
		Get(addr)

	return resp, err
}

func (g *Gitee) put(addr string, data map[string]string) (*resty.Response, error) {
	data["access_token"] = g.Token

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(mapToJSONString(data)).
		Put(addr)

	return resp, err
}

func (g *Gitee) post(addr string, data map[string]string) (*resty.Response, error) {
	data["access_token"] = g.Token

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(mapToJSONString(data)).
		Post(addr)

	return resp, err
}

func (g *Gitee) info(resp *resty.Response, err error) {
	// Explore response object
	fmt.Println("Response Info:")
	fmt.Println("  Error      :", err)
	fmt.Println("  Status Code:", resp.StatusCode())
	fmt.Println("  Status     :", resp.Status())
	fmt.Println("  Body       :\n", resp)
	fmt.Println()

	// Explore trace info
	// fmt.Println("Request Trace Info:")
	// ti := resp.Request.TraceInfo()
	// fmt.Println("  DNSLookup     :", ti.DNSLookup)
	// fmt.Println("  ConnTime      :", ti.ConnTime)
	// fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
	// fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
	// fmt.Println("  ServerTime    :", ti.ServerTime)
	// fmt.Println("  ResponseTime  :", ti.ResponseTime)
	// fmt.Println("  TotalTime     :", ti.TotalTime)
	// fmt.Println("  IsConnReused  :", ti.IsConnReused)
	// fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
	// fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
}
