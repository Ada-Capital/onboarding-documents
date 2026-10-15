# Prompt templates
# ---------------------------------------------------------------------------

SYSTEM_PROMPT = """Você é um analista de documentos de compliance de uma fintech brasileira regulada.
Você recebe imagens de um documento enviado por um candidato a funcionário e os dados que ele declarou no formulário.
Analise apenas o que está visível no documento. Não avalie autenticidade, fraude ou qualidade de foto.
Responda somente com um objeto JSON válido, sem markdown e sem texto fora do JSON."""

CNH_TYPE_TEMPLATE = """Modelo CNH (Carteira Nacional de Habilitação, Brasil):
- Título "CARTEIRA NACIONAL DE HABILITAÇÃO" e brasão da República; emitida pelo DETRAN de um estado.
- Campos: NOME; DOC. IDENTIDADE / ÓRG. EMISSOR / UF; CPF; DATA NASCIMENTO; FILIAÇÃO;
  Nº REGISTRO (11 dígitos); VALIDADE; 1ª HABILITAÇÃO; CAT. HAB. (A, B, AB, C, D, E);
  LOCAL e DATA DE EMISSÃO; foto do titular; assinatura do titular.
- Modelos a partir de 2017 têm QR code; o modelo 2022 pode trazer MRZ e o texto "DRIVER LICENSE".
- Pode ser o documento físico fotografado ou o PDF da CNH Digital (app Carteira Digital de Trânsito)."""

RG_TYPE_TEMPLATE = """Modelo RG (Carteira de Identidade, Brasil) — frente e verso:
- Frente: "REPÚBLICA FEDERATIVA DO BRASIL", estado emissor, órgão expedidor (SSP, IIRGD, PC, DETRAN...),
  foto do titular, impressão digital e assinatura.
- Verso: REGISTRO GERAL (número do RG); DATA DE EXPEDIÇÃO; NOME; FILIAÇÃO; NATURALIDADE;
  DATA DE NASCIMENTO; DOC. ORIGEM (certidão); CPF (opcional em modelos antigos).
- Modelo novo CIN (Carteira de Identidade Nacional, a partir de 2022): o número do documento é o próprio CPF,
  traz QR code e MRZ, pode ser físico ou o PDF do app gov.br."""

PROOF_OF_ADDRESS_EXAMPLE = """Modelo de comprovante de residência aceito:
- Conta de consumo (água, luz, gás, telefone fixo/celular, internet/TV), fatura de cartão de crédito,
  extrato ou correspondência bancária, boleto de condomínio/IPTU, contrato de aluguel ou declaração
  de residência assinada.
- Deve conter: nome do titular; endereço completo (logradouro, número, complemento, bairro, cidade,
  estado e CEP); nome/logotipo do emissor; data de emissão ou referência.
- Data de emissão de no máximo 90 dias em relação à data da análise.
- Titular deve ser o próprio candidato. Comprovante em nome de terceiro (pai, mãe, cônjuge) só é
  aceito com o mesmo sobrenome ou vínculo evidente, e deve ser sinalizado na justificativa."""

IDENTITY_PROMPT = """Analise o documento anexado e responda se ele é do tipo RG ou CNH e se os campos informados estão coerentes.

{{CNH_TYPE_TEMPLATE}}

{{RG_TYPE_TEMPLATE}}

Documento anexado: {{DOC_FILE}}

Campos informados pelo candidato no formulário:
{{FILLED_FIELDS}}

Regras:
- Compare nome completo, data de nascimento, CPF e número do documento (RG/CNH) do documento com os campos informados.
  Tolere diferenças de acentuação, abreviação e formatação; divergência de pessoa (nome E data de nascimento ou CPF diferentes) invalida.
- Se o candidato não for brasileiro (nationality diferente de BRA) o documento pode ser passaporte ou identidade estrangeira:
  classifique como "OTHER" e valide nome, data de nascimento e número do documento contra os campos informados.
- Documento ilegível, incompleto (RG sem o verso) ou que não seja um documento de identificação é inválido.
- Não avalie autenticidade, fraude ou qualidade da foto; avalie apenas tipo e coerência dos campos.
- Data da análise: {{ANALYSIS_DATE}}.

Responda apenas em JSON:
{"is_valid": true | false, "document_type": "RG" | "CNH" | "OTHER", "reason": "justificativa objetiva em português, 1 a 3 frases"}"""

ADDRESS_PROMPT = """Analise o documento anexado e responda se ele é um comprovante de residência válido.

{{DOC_FILE_EXAMPLE}}

Documento anexado: {{DOC_FILE}}

Campos informados pelo candidato no formulário:
{{FILLED_FIELDS}}

Regras:
- O nome do titular deve corresponder ao nome informado (tolere abreviações e acentuação).
- O endereço do documento deve corresponder ao endereço informado (logradouro, número, cidade, estado e CEP);
  abreviações e sinônimos são aceitos; rua ou cidade diferentes invalidam.
- Data de emissão de no máximo 90 dias em relação à data da análise: {{ANALYSIS_DATE}}. Documento sem data é inválido.
- Documento ilegível ou que não seja um comprovante de residência é inválido.
- Não avalie autenticidade ou fraude; avalie apenas tipo, titularidade, endereço e data.

Responda apenas em JSON:
{"is_valid": true | false, "reason": "justificativa objetiva em português, 1 a 3 frases"}"""

