import XGStudio from "../contracts/XGStudio.cdc"

pub fun main():{UInt64:XGStudio.Template}  {
    return XGStudio.getAllTemplates()
}
