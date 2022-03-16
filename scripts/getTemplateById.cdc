import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"


pub fun main(templateId: UInt64): TroonAtomicStandard.Template {
    return TroonAtomicStandard.getTemplateById(templateId: templateId)
}