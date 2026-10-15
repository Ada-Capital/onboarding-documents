# Prompts de análise de documento (origem: 4Platforms)

Cópia literal de `Ada-Capital/BRLAForPlatformsLambdas@80a4452` (branch `production`), salvo onde a tabela indica outra branch. Prompts que moram dentro do código foram recortados do `app.py`/`llm.py` (as variáveis `{...}` são f-strings do original). Não editar aqui sem anotar a diferença: a fonte de verdade continua sendo o repo de origem.

| Arquivo | Origem | Serve para (A4) |
|---|---|---|
| `kyb-subs/SYSTEM-incorporation.txt` | `lambdas/kyb-subs-analyzer/prompts/` | Contrato social / estatuto |
| `kyb-subs/SYSTEM-ubo-identity.txt` | idem | Identificação PF (UBO / representante) |
| `kyb-subs/SYSTEM-proof-of-address.txt` | idem | Comprovante de endereço |
| `kyb-subs/SYSTEM-financial-capacity.txt` | idem | Capacidade financeira |
| `kyb-subs/SYSTEM-company-tax.txt` | idem | Cartão CNPJ / tax ID (não está na lista do plano) |
| `kyb-subs/SYSTEM-limit-analysis.txt` | idem | Limite recomendado a partir de extratos, sob demanda |
| `kyb-subs/resources/*.json` | `lambdas/_kyb_subs_core/kyb_subs/resources/` | Envelope de saída, tabela de decisão e parâmetros usados pelos prompts acima |
| `kyb-subs/PROMPT-VERSIONS.txt` | `docs/features/kyb-subs/` | Versão de cada prompt |
| `document-parser-kyc/prompt.py` + `models.py` | `lambdas/document_parser_kyc/document_parser/` | Identidade BR (RG, CIN, CNH, passaporte, CPF): tipo, campos, qualidade (A4; também separa ID de outros docs no A3) |
| `financial-capacity-kyc/classifier_prompt.py` | `lambdas/financial_capacity_kyc/financial_capacity/llm.py` | Classificador de tipo de documento financeiro (base para o classificador do A3) |
| `financial-capacity-kyc/prompt.py` | `lambdas/financial_capacity_kyc/evals/` | Prompt do eval do classificador |
| `proof-of-address-validator/poa_extraction.py` | branch `fork/poa-analyze`, `lambdas/validate-proof-of-address/prompts/` | Comprovante de endereço com evidência por campo (alternativa ao `SYSTEM-proof-of-address`) |
| `sumsub-documents-analyser/prompts.py` | `lambdas/sumsub-documents-analyser/app.py` | Documento × dados declarados (ID e comprovante, KYE) |
| `due-diligence-kyc/prompt.py` | `lambdas/due_diligence_kyc/app.py` | Pesquisa web de risco 1–5 por nome/CPF (background check, A7) |
| `company-summary/prompt.py` | `lambdas/company_summary/app.py` | Resumo da empresa a partir do site (A2 pré-preenchimento, A7) |
| `rfi-taxonomy/rfiCategories.go` | `BRLAForPlatforms`, branch `feature/internal-kyb-endpoint`, `Lib/pkg/bridgeRfi/` | 20 categorias de RFI de parceiro com mensagem ao cliente em pt-BR (A5). Classificador por regex, sem LLM |

Os códigos de tipo de documento que o backend aceita (o A3 deve classificar nestes): `ID`, `DRIVERS-LICENSE`, `PASSPORT`, `RESIDENCE-PERMIT`, `PROOF-OF-ADDRESS`, `PROOF-OF-ADDRESS-COMPANY`, `PROOF-OF-FINANCIAL-CAPACITY`, `PROOF-OF-REVENUE`, `CERTIFICATE-OF-INCORPORATION`, `ARTICLES-OF-INCORPORATION`, `COMPANY-TAX-IDENTIFICATION-DOCUMENT`, `OWNERSHIP-CHART`, `PROOF-OF-SIGNING-AUTHORITY`, `SHAREHOLDER-REGISTER`, `DIRECTORS-REGISTER`, `INVOICE`, `ADDITIONAL-ATTACHMENT` (`Lib/entities/uploadedDocument.go`).

Sem prompt pronto: source of funds, org chart e poderes dos signatários (o `incorporation` extrai só nome + cargo dos administradores).

## Como chamar

- System = o `SYSTEM-*.txt` inteiro. Não há placeholders.
- User = envelope JSON (`analysis_date`, `slot`, `declared`, `document`, `params`, `versions`, `context`) + as páginas como imagem. Formato em `resources/envelopes.json`; `params` vem de `resources/params-defaults.json` (o prompt não usa valor próprio).
- `financial-capacity` e `limit-analysis` precisam de `params.fx_rates`; sem isso devolvem `derived.fx_requests`.
- Produção usa `temperature=0`, PDF → PNG 150 dpi, até 20 páginas.

## Saída → status da exigência

O modelo nunca decide; `resources/decision-table.json` decide pelos `findings[].code`:

| Decisão | `requirements.status` |
|---|---|
| APPROVE (sem finding ou só `INFO_NOTE`) | `ok` |
| REJECT | `not_ok` (`DOC_OUTDATED` → `expired`) |
| REVIEW (finding REVIEW ou `confidence=low`) | fica para humano |

`findings[].rfi_suggestion` já traz o texto do que pedir ao cliente (entrada do A5).
