package bridgeRfi

import (
	"html"
	"regexp"
	"sort"
	"strings"
	"unicode"

	"BRLAForPlatforms/Lib/entities"
)

// Artifact kinds an RFI item can demand from us.
const (
	ArtifactDocument = "DOCUMENT"
	ArtifactAnswer   = "ANSWER"
	ArtifactLink     = "LINK"
	ArtifactForm     = "FORM"
)

// maxRfiBodyRunes caps a pasted body before any regex work, so a pathological paste cannot
// turn the classifier into a CPU sink.
const maxRfiBodyRunes = 200000

// minSpanLetters is the shortest span we bother reporting back as unclassified — below that a
// span is a salutation, a signature line or a date, never a request.
const minSpanLetters = 12

// RfiCategory is one entry of the operator RFI taxonomy: a request kind we have actually seen,
// with the neutral pt-BR text we may show a customer and the artifact it demands.
type RfiCategory struct {
	ID              int
	Slug            string
	Title           string
	CustomerMessage string
	Artifact        string
	DocTypeHint     string
	Actor           string
	Severity        string
	Anticipatable   bool

	matchers []*regexp.Regexp
}

// frag turns a literal e-mail fragment into a case-insensitive matcher whose runs of whitespace
// match any whitespace, so line wrapping in the pasted e-mail does not break the match.
func frag(literal string) *regexp.Regexp {
	quoted := regexp.QuoteMeta(strings.TrimSpace(literal))
	quoted = regexp.MustCompile(`(?:\\ |\s)+`).ReplaceAllString(quoted, `\s+`)
	return regexp.MustCompile(`(?i)` + quoted)
}

func frags(literals ...string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, 0, len(literals))
	for _, l := range literals {
		out = append(out, frag(l))
	}
	return out
}

// categories is the 20-entry taxonomy. ID is the operator sheet's own numbering and is the stable
// id ops refer to; Slug is the stable machine id. Messages are written from the sheet's description
// column and never name the partner, an external questionnaire product or an external URL.
var categories = []RfiCategory{
	{
		ID:              1,
		Slug:            "nature-of-business-evidence",
		Title:           "Comprovação da natureza do negócio",
		CustomerMessage: "Envie notas fiscais, contratos ou acordos de serviço firmados entre a empresa e seus clientes, que comprovem a atividade declarada e os produtos ou serviços prestados. Os documentos precisam trazer os dados completos da empresa emitente.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.INVOICE),
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"invoices, contracts, and/or service agreements",
			"sales invoices, contracts",
			"substantiate their nature of business",
			"to supplement their website",
		),
	},
	{
		ID:              2,
		Slug:            "business-website",
		Title:           "Site da empresa ou como capta clientes",
		CustomerMessage: "Informe o link do site institucional da empresa. Se a empresa não tiver site, ou se o site informado estiver fora do ar, descreva como ela se divulga ao mercado e capta clientes.",
		Artifact:        ArtifactLink,
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"link to the business",
			"how the applying entity markets themselves",
			"markets themselves to generate business",
			"provided website is not working",
		),
	},
	{
		ID:              3,
		Slug:            "operating-address-proof",
		Title:           "Comprovação do endereço operacional",
		CustomerMessage: "O endereço informado foi identificado como endereço de agente registrado ou endereço virtual. Envie o contrato de locação vigente ou uma conta de consumo emitida nos últimos 90 dias. Se a operação for remota, é possível usar o endereço residencial de um sócio controlador — nesse caso o comprovante deve estar no nome dele. Extrato bancário não é aceito como comprovante de endereço.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.PROOFOFADDRESSCOMPANY),
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"known registered agent or virtual address",
			"operating address provided",
			"lease agreement or utility bill",
			"utility bill < 90 days old",
		),
	},
	{
		ID:              4,
		Slug:            "owner-kyc-photo-id",
		Title:           "KYC e documento de sócio ou controlador com 25% ou mais",
		CustomerMessage: "É necessário concluir o cadastro de identificação e enviar a foto do documento de identidade de cada pessoa com 25% ou mais de participação, além dos controladores. Se a participação já estiver abaixo de 25%, envie o ato societário atualizado que comprove a alteração. Documentos vencidos não são aceitos.",
		Artifact:        ArtifactForm,
		DocTypeHint:     string(entities.ID),
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"KYC information and photo ID",
			"25%+ owner or control person",
			"no longer affiliated",
			"share a current filing",
		),
	},
	{
		ID:              5,
		Slug:            "licensure-proof",
		Title:           "Prova de licenciamento ou de isenção",
		CustomerMessage: "Envie a comprovação de que a empresa é licenciada ou registrada para exercer as atividades declaradas. Se nenhuma licença for exigida, envie uma explicação detalhada da isenção, indicando a base legal ou regulatória que permite operar sem licença. Comprovante de inscrição fiscal não substitui a licença da atividade.",
		Artifact:        ArtifactDocument,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"proof of licensure",
			"explanation of their exemption",
			"registered to provide their stated activities",
			"regulatory basis for operating without a license",
		),
	},
	{
		ID:              6,
		Slug:            "volume-spike-explanation",
		Title:           "Explicação de volume acima do declarado",
		CustomerMessage: "O volume transacionado excedeu o volume declarado no cadastro. Explique a discrepância, informe se foi um evento pontual e apresente a projeção de volume mensal para os próximos meses. Responda cada item com pelo menos uma ou duas frases, detalhando o propósito das transações e o uso dos recursos.",
		Artifact:        ArtifactAnswer,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"substantially exceeded the volume declared",
			"spike in activity observed",
			"transaction activity observed on the account",
		),
	},
	{
		ID:              7,
		Slug:            "source-of-wealth-evidence",
		Title:           "Evidência documental de origem do patrimônio",
		CustomerMessage: "Envie documentos que comprovem como o patrimônio foi acumulado, por exemplo contratos de trabalho, holerites, demonstrações financeiras, carteira de investimentos ou documentos de herança.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.PROOFOFFINANCIALCAPACITY),
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"documentary evidence substantiating your source of wealth",
			"employment contracts, payslips",
			"inheritance documents",
		),
	},
	{
		ID:              8,
		Slug:            "source-of-funds-questionnaire",
		Title:           "Questionário de origem dos recursos",
		CustomerMessage: "É necessário responder ao questionário de origem dos recursos. Ele é exigido quando o endereço de residência do titular está em país de alto risco ou quando o titular tem 60 anos ou mais.",
		Artifact:        ArtifactForm,
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"Source of Funds Questionnaire",
			"SOFQ",
			"address belongs to a high-risk country",
			"60 years of age or older",
		),
	},
	{
		ID:              9,
		Slug:            "counterparty-purpose",
		Title:           "Natureza e propósito das transações com contraparte",
		CustomerMessage: "Explique a relação comercial com a contraparte indicada: quais bens ou serviços estão envolvidos, qual o papel de eventuais intermediários e qual o propósito das transações. Anexe documentos de suporte, como recibos, notas fiscais, confirmações de pedido, comprovantes de entrega ou correspondência com a contraparte.",
		Artifact:        ArtifactAnswer,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"nature and purpose of your transactions",
			"role of any intermediaries",
			"businesses you have interacted with",
		),
	},
	{
		ID:              10,
		Slug:            "ownership-structure",
		Title:           "Composição societária e percentuais de participação",
		CustomerMessage: "Envie um documento formal que liste os sócios ou acionistas da empresa e os respectivos percentuais de participação. Quando houver pessoa jurídica entre os sócios relevantes, a mesma comprovação é exigida para essas entidades intermediárias.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.SHAREHOLDERREGISTER),
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"shareholders of the company and their ownership percentages",
			"ownership structure for both of these entities",
			"provide ownership structure",
		),
	},
	{
		ID:              11,
		Slug:            "third-party-account-use",
		Title:           "Uso da conta para terceiros ou clientes finais",
		CustomerMessage: "Confirme se a conta será usada para investir, transferir ou negociar recursos por conta de clientes ou de terceiros. Se não for esse o caso, detalhe como a empresa pretende usar a conta.",
		Artifact:        ArtifactAnswer,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"on behalf of customers or 3rd-parties",
			"used to invest, transfer, or trade funds",
			"how the entity will use their account",
		),
	},
	{
		ID:              12,
		Slug:            "source-of-funds-narrative",
		Title:           "Explicação da origem dos recursos e de incoerência cadastral",
		CustomerMessage: "Explique a origem dos recursos e esclareça a divergência entre as informações declaradas no cadastro, por exemplo entre a ocupação informada e a fonte de renda. Quando a receita vier de doações ou fomento, confirme se há investidores externos envolvidos.",
		Artifact:        ArtifactAnswer,
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"Kindly clarify this discrepancy",
			"source of funds is indicated as",
			"employment status is listed as",
			"includes any external investors",
		),
	},
	{
		ID:              13,
		Slug:            "source-of-funds-documents",
		Title:           "Comprovação documental da origem dos recursos",
		CustomerMessage: "Envie um documento recente entre as opções aceitas: demonstrações financeiras auditadas, extrato bancário ou de corretora, ou declaração de imposto. Para pessoa física, o extrato bancário, de corretora ou de cripto deve ter menos de 90 dias. Demonstrações financeiras precisam ser auditadas e assinadas.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.PROOFOFFINANCIALCAPACITY),
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"confirming source of funds",
			"Audited Financials",
			"Bank/Brokerage account statement",
			"Tax return",
		),
	},
	{
		ID:              14,
		Slug:            "business-formation-document",
		Title:           "Ato constitutivo da empresa",
		CustomerMessage: "Envie o ato constitutivo da empresa, como o certificado de constituição ou o contrato ou estatuto social, atendendo aos requisitos formais de conteúdo e assinatura.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.CERTIFICATEOFINCORPORATION),
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"business formation document",
			"certificate of formation",
			"articles of incorporation",
		),
	},
	{
		ID:              15,
		Slug:            "entity-tax-id-confirmation",
		Title:           "Confirmação do número de identificação fiscal da empresa",
		CustomerMessage: "Envie um documento que confirme o número de identificação fiscal da empresa emitido pela autoridade fiscal do país de constituição. Para empresas constituídas nos Estados Unidos, a carta IRS 147(c) é o documento esperado.",
		Artifact:        ArtifactDocument,
		DocTypeHint:     string(entities.COMPANYTAXIDENTIFICATIONDOCUMENT),
		Actor:           actorClient,
		Severity:        severityBlocking,
		Anticipatable:   true,
		matchers: frags(
			"EIN/Tax ID of the entity",
			"IRS 147(c)",
			"147(c) letter",
		),
	},
	{
		ID:              16,
		Slug:            "third-party-site-affiliation",
		Title:           "Vínculo com empresa ou site de terceiro",
		CustomerMessage: "O site informado está associado a outra razão social. Comprove a relação entre a empresa solicitante e essa outra empresa, por meio de nome fantasia registrado, acordo de serviço ou organograma. Como alternativa, comprove a titularidade do domínio com o registro do domínio, uma fatura recente de hospedagem ou uma captura do painel administrativo do site.",
		Artifact:        ArtifactDocument,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"the company affiliated with the website",
			"entity's relationship with",
			"domain registration details",
			"hosting account invoice",
			"website administration panel",
		),
	},
	{
		ID:              17,
		Slug:            "anti-fraud-questions",
		Title:           "Perguntas antifraude ao titular",
		CustomerMessage: "O próprio titular precisa responder: se alguma pessoa, inclusive familiar, pediu que ele abrisse esta conta; se foi convencido a enviar transferências ou abrir contas por causa de uma suposta oportunidade de negócio; e se concedeu a terceiros acesso ou controle das suas finanças.",
		Artifact:        ArtifactAnswer,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"ask you to set up this account",
			"business opportunity",
			"convinced to send money through wire transfers",
			"control of your finances",
		),
	},
	{
		ID:              18,
		Slug:            "sector-import-logistics-documents",
		Title:           "Documentação setorial de importação e logística",
		CustomerMessage: "Para a atividade de importação, envie os contratos de locação dos armazéns onde os produtos são guardados, os conhecimentos de embarque ou contratos de frete que comprovem os embarques, e os contratos assinados com fornecedores, distribuidores ou parceiros logísticos. Se algum desses documentos não existir, explique o motivo.",
		Artifact:        ArtifactDocument,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"Warehouse leases",
			"Bills of lading",
			"shipping contracts",
			"freight contracts",
		),
	},
	{
		ID:              19,
		Slug:            "elaborate-nature-of-business",
		Title:           "Detalhamento da natureza do negócio",
		CustomerMessage: "A descrição da atividade da empresa foi considerada vaga. Detalhe o que a empresa faz: produtos ou serviços oferecidos, perfil dos clientes, como a receita é gerada e em quais países opera.",
		Artifact:        ArtifactAnswer,
		Actor:           actorClient,
		Severity:        severityBlocking,
		matchers: frags(
			"elaborate and expand upon the nature of the business",
			"expand upon the nature of the business",
		),
	},
	{
		ID:              20,
		Slug:            "past-due-follow-up",
		Title:           "Cobrança de solicitação em atraso",
		CustomerMessage: "Há uma solicitação de informações anterior ainda sem resposta, com novo prazo definido. Enquanto ela não for atendida, o processamento das transações pode ser afetado.",
		Artifact:        ArtifactAnswer,
		Actor:           actorOps,
		Severity:        severityInformational,
		matchers: frags(
			"past due request",
			"follow up regarding our past due",
			"continue processing the user's transactions",
		),
	},
}

// Categories returns the taxonomy ordered by ID. The copies carry no usable matcher slice for
// callers to mutate through — matchers stay an implementation detail of this package.
func Categories() []RfiCategory {
	out := make([]RfiCategory, len(categories))
	copy(out, categories)
	for i := range out {
		out[i].matchers = nil
	}
	return out
}

func CategoryBySlug(slug string) (RfiCategory, bool) {
	for _, c := range categories {
		if c.Slug == slug {
			c.matchers = nil
			return c, true
		}
	}
	return RfiCategory{}, false
}

var (
	htmlTagRe    = regexp.MustCompile(`(?s)<[^>]{0,1000}>`)
	quotedLineRe = regexp.MustCompile(`(?im)^\s*(?:>|On .{0,200}wrote:|Em .{0,200}escreveu:).*$`)
	bulletRe     = regexp.MustCompile(`^\s*(?:[-*•·–]+|\d{1,2}[.)])\s*`)
	sentenceRe   = regexp.MustCompile(`[.!?]+\s+`)
	scriptRe     = regexp.MustCompile(`(?is)<(?:script|style)\b.{0,1000}?</(?:script|style)>`)
	blockTagRe   = regexp.MustCompile(`(?i)<\s*br\s*/?\s*>|</\s*(?:p|div|li|tr|h[1-6])\s*>`)
	hSpaceRe     = regexp.MustCompile(`[ \t\x{00a0}]+`)
)

// ClassifyRfiText classifies one pasted RFI e-mail body. One e-mail commonly carries several
// requests, so every matching category is returned, deduped and ordered by ID.
//
// Matching runs over the whole normalized body (a request may wrap across lines). The unclassified
// report is then computed per span — the body is split on blank lines, line breaks, bullet markers
// and sentence boundaries — and a span is reported verbatim only when it overlaps no range that a
// matcher actually claimed, so the operator sees a real miss and never sees a request that was
// classified fine but happened to wrap mid-phrase.
func ClassifyRfiText(body string) (matched []RfiCategory, unclassified []string) {
	matched = []RfiCategory{}
	unclassified = []string{}

	normalized := normalizeBody(body)
	if normalized == "" {
		return matched, unclassified
	}

	seen := map[int]struct{}{}
	for _, c := range categories {
		if matchesAny(c.matchers, normalized) {
			if _, dup := seen[c.ID]; dup {
				continue
			}
			seen[c.ID] = struct{}{}
			hit := c
			hit.matchers = nil
			matched = append(matched, hit)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].ID < matched[j].ID })

	covered := coveredRanges(normalized)
	cursor := 0
	seenSpan := map[string]struct{}{}
	for _, span := range splitSpans(normalized) {
		// Spans are substrings of normalized, so a moving cursor locates each one in order.
		start := strings.Index(normalized[cursor:], span)
		if start >= 0 {
			start += cursor
			cursor = start + len(span)
			if overlapsCovered(covered, start, start+len(span)) {
				continue
			}
		} else if categoryMatchesSpan(span) {
			continue
		}
		if _, dup := seenSpan[span]; dup {
			continue
		}
		seenSpan[span] = struct{}{}
		unclassified = append(unclassified, span)
	}

	return matched, unclassified
}

// coveredRanges collects every byte range of the normalized body claimed by some category matcher,
// so a span overlapping one is known-classified rather than a miss.
func coveredRanges(normalized string) [][]int {
	ranges := make([][]int, 0, 8)
	for _, c := range categories {
		for _, m := range c.matchers {
			ranges = append(ranges, m.FindAllStringIndex(normalized, -1)...)
		}
	}
	return ranges
}

func overlapsCovered(covered [][]int, start, end int) bool {
	for _, r := range covered {
		if start < r[1] && r[0] < end {
			return true
		}
	}
	return false
}

func matchesAny(matchers []*regexp.Regexp, text string) bool {
	for _, m := range matchers {
		if m.MatchString(text) {
			return true
		}
	}
	return false
}

func categoryMatchesSpan(span string) bool {
	for _, c := range categories {
		if matchesAny(c.matchers, span) {
			return true
		}
	}
	return false
}

// normalizeBody strips HTML markup and quoted reply blocks, unescapes entities and collapses
// horizontal whitespace, keeping line structure for span splitting.
func normalizeBody(body string) string {
	runes := []rune(body)
	if len(runes) > maxRfiBodyRunes {
		body = string(runes[:maxRfiBodyRunes])
	}

	body = strings.ReplaceAll(body, "\r\n", "\n")
	body = scriptRe.ReplaceAllString(body, " ")
	body = blockTagRe.ReplaceAllString(body, "\n")
	body = htmlTagRe.ReplaceAllString(body, " ")
	body = html.UnescapeString(body)
	body = quotedLineRe.ReplaceAllString(body, "")

	lines := strings.Split(body, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(hSpaceRe.ReplaceAllString(line, " "))
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func splitSpans(normalized string) []string {
	spans := make([]string, 0, 16)
	for _, line := range strings.Split(normalized, "\n") {
		line = bulletRe.ReplaceAllString(line, "")
		for _, sentence := range sentenceRe.Split(line, -1) {
			sentence = strings.TrimSpace(strings.Trim(strings.TrimSpace(sentence), "-*•·–"))
			if countLetters(sentence) < minSpanLetters {
				continue
			}
			spans = append(spans, sentence)
		}
	}
	return spans
}

func countLetters(s string) int {
	n := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			n++
		}
	}
	return n
}
