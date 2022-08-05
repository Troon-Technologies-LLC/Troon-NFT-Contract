import XGStudio from "../contracts/XGStudio.cdc"


pub fun main() : Int {
  return  XGStudio.getAllTemplates().length
}