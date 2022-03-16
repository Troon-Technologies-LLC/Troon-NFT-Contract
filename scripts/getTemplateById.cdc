import TroonAtomicStandardContract from "../contracts/TroonAtomicStandardContract.cdc"


pub fun main(templateId: UInt64): TroonAtomicStandardContract.Template {
    return TroonAtomicStandardContract.getTemplateById(templateId: templateId)
}