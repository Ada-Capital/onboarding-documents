"""Prompt de classificação para a camada LLM (modelo leve).

Só classifica e justifica em uma frase — os números continuam saindo da extração
determinística. Modelo não inventa valor.
"""

TAXONOMIA = """balance_sheet: balanço patrimonial (ativo, passivo e patrimônio líquido, com data de referência)
income_statement: DRE / demonstração do resultado do exercício
trial_balance: balancete de verificação (contas com débito/crédito/saldo)
opening_balance: balanço de abertura de empresa nova
interim_statements: demonstrações intermediárias de período parcial
audited_statements: demonstrações acompanhadas de parecer de auditor independente
revenue_report: relação/relatório de faturamento emitido pela contabilidade, com meses realizados
bank_statement: extrato de conta (lista de lançamentos/movimentações com datas e valores), de banco ou instituição de pagamento, em qualquer idioma
projection_or_business_plan: plano de negócios, projeção ou meta de receita futura
incorporation_doc: contrato social, alteração contratual, requerimento de empresário
utility_bill_or_invoice: fatura, boleto, conta de consumo (luz, água, telefone) ou fatura de cartão
illegible: sem texto aproveitável — scan sem OCR ou texto corrompido/embaralhado
unrecognized: documento legível que não é nenhum dos acima (ex.: declaração de que a conta existe, sem movimentação)"""

REGRAS = """Desempates que importam:
- extrato tem LANÇAMENTOS (linhas de data + histórico + valor). Documento de banco que só declara que a conta existe, sem lançamentos, é unrecognized.
- fatura/boleto/conta de consumo é utility_bill_or_invoice mesmo quando emitida por banco ou operadora e mesmo que traga lançamentos de cartão.
- se o texto está embaralhado a ponto de não dar pra ler os campos, é illegible, não unrecognized.
- documento em inglês ou espanhol segue a mesma taxonomia."""

SYSTEM = ("Você classifica documentos de KYB para uma instituição financeira brasileira. "
          "Responda só com JSON válido, sem markdown.")


def build(texto: str, max_chars: int = 2500) -> str:
    return f"""Classifique o documento abaixo em exatamente um tipo da taxonomia.

TAXONOMIA
{TAXONOMIA}

{REGRAS}

Responda no formato JSON: {{"doc_type": "<tipo>", "confianca": <0-100>, "justificativa": "<uma frase, citando o que no texto sustenta a escolha>"}}

TEXTO DO DOCUMENTO (primeiros {max_chars} caracteres):
---
{texto[:max_chars]}
---"""
