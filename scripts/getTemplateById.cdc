import NFTContractV01 from "../contracts/NFTContractV01.cdc"


pub fun main(templateId: UInt64): NFTContractV01.Template {
    return NFTContractV01.getTemplateById(templateId: templateId)
}