from typing import Literal, Optional
from pydantic import BaseModel, Field


class DocumentFields(BaseModel):
    """Union of all possible fields across all document types.
    Gemini fills only the fields applicable to the identified document type.
    """
    # Common
    full_name: Optional[str] = None
    date_of_birth: Optional[str] = Field(None, description="Format YYYY-MM-DD")
    photo_present: Optional[bool] = None

    # RG
    rg_number: Optional[str] = None
    father_name: Optional[str] = None
    mother_name: Optional[str] = None
    place_of_birth: Optional[str] = Field(None, description="City/State")
    issuing_authority: Optional[str] = Field(None, description="E.g. SSP/SP, IFP/RJ, DETRAN/SP")
    issue_date: Optional[str] = Field(None, description="Format YYYY-MM-DD")
    cpf: Optional[str] = Field(None, description="CPF if visible on RG")

    # CIN
    social_name: Optional[str] = Field(None, description="If present on CIN")
    gender: Optional[str] = None
    nationality: Optional[str] = None
    national_registration_number: Optional[str] = None
    cpf_number: Optional[str] = None
    expiration_date: Optional[str] = Field(None, description="Format YYYY-MM-DD")
    mrz: Optional[str] = Field(None, description="Machine readable zone text")
    qr_code_present: Optional[bool] = None

    # CNH
    registration_number: Optional[str] = None
    license_number: Optional[str] = None
    category: Optional[str] = Field(None, description="A, B, AB, C, D, E, etc.")
    first_license_date: Optional[str] = Field(None, description="Format YYYY-MM-DD")
    issue_location: Optional[str] = Field(None, description="City/State")
    document_notes: Optional[str] = Field(None, description="CNH Obs field, e.g. EAR")

    # Passport
    passport_number: Optional[str] = None
    passport_issuing_authority: Optional[str] = None
    mrz_line1: Optional[str] = None
    mrz_line2: Optional[str] = None

    # CPF
    registration_date: Optional[str] = Field(None, description="Format YYYY-MM-DD, if visible")


class QualityAssessment(BaseModel):
    image_quality: Literal["GOOD", "FAIR", "POOR"]
    document_readable: bool
    partially_readable_fields: list[str] = Field(default_factory=list)


class DocumentParserResponse(BaseModel):
    document_type: Literal["RG", "CIN", "CNH", "PASSPORT", "CPF", "UNKNOWN"]
    identification_confidence: float = Field(ge=0.0, le=1.0)
    document_side: Literal["FRONT", "BACK", "BOTH"]
    fields: DocumentFields
    quality: QualityAssessment
    observations: Optional[str] = None
