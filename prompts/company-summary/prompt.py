        prompt = f"""
You will receive a URL containing the website of a company. Visit the URL, analyze all publicly accessible information, and generate a clear, professional, and factual summary of what the company does.

Instructions:

1. Review relevant sections such as About, Services, Products, Solutions, Mission, and any other page describing the company.

2. Summarize the company's core activities, business model, services or products offered, target market, industry sector, and strategic purpose.

3. The summary must always be written in English, even if the company is foreign or the website content is in a different language.

4. Write the summary in the same tone, structure, and level of detail as the following example:

Example Style:
"EZCAPITAL INVESTIMENTOS LTDA is an investment and holding company established to manage its own capital and hold equity interests in other businesses. Its primary activities include acquiring, maintaining, and administering ownership stakes (quotas/shares) in companies in Brazil and abroad, acting as a holding vehicle for investment, control, and strategic management purposes. The company may also carry out financial operations related to its investment activity, such as capital allocation, corporate structuring, and governance support for its investee companies, with a focus on capital appreciation and portfolio development."

5. The summary must be 5–8 sentences, objective, neutral, and based solely on what is present on the website—no assumptions.

6. DO NOT MENTION URL OR WEBSITE LINK IN THE SUMMARY.

Input URL: {url}
        """
