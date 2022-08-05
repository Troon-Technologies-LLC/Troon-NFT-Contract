import XGStudio from "../contracts/XGStudio.cdc"

pub fun main(templateId: UInt64): XGStudio.Template {
    return XGStudio.getTemplateById(templateId: templateId)
}