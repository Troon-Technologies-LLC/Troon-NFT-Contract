import TroonAtomicStandard from "../contracts/TroonAtomicStandard.cdc"

pub fun main():{UInt64:TroonAtomicStandard.Template}  {
    return TroonAtomicStandard.getAllTemplates()
}
