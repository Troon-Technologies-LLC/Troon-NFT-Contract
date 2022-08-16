import NFTContract from 0x179b6b1cb6755e31
pub fun main(templateId:UInt64): AnyStruct{
    return NFTContract.getTemplateById(templateId: templateId)
}