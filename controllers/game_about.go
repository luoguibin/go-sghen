package controllers

//  order type key          value   description
const (
    OTypeMsg            =   10000   // message order
    OTypeMsgPerson      =   11000   // msg from client to client
    OTypeMsgNear        =   12000   // msg from client to near clients
    OTypeMsgGroup       =   13000   // msg from client to a group
    OTypeMsgAll         =   14000   // msg from client to all
    OTypeMsgSystem      =   15000   // msg from system to client ro clients
    OTypeMsgSystemMaintenance   =   15001
    OTypeMsgSystemPerson        =   15002

    OTypeSkill          =   20000   // skill order
    OTypeSkillSingle    =   21000   // single target skill with damage
    OTypeSkillSingleK   =   22000   // single target skill with control
    OTypeSkillNear      =   25000   // near targets skill with damage
    OTypeSkillNearK     =   27000   // near targets skill with control
    
    OTypeAction         =   30000   // action order
    OTypeActionDrug     =   31000   // take some drug
    OTypeActionMove     =   32000   // client actions that relative to move 

    OTypeData           =   40000   // game data signal order
    OTypeDataPerson     =   40001   // update personal game data
    OTypeDataAll        =   40002   // update clients game data
)  

//  id type key             value   description
const (
    ITypePerson         =   10000   // client id group
    ITypeGroup          =   20000   // clients id group
    ITypeSystem         =   30000   // system id group

    IDSYSTEM    int64   =   30000   // system id
    IDPERSON    int64   =   10000   // client id
)

const (
    GStatusLogin 	    =	100
	GStatusInGame	    =   101
	GStatusLogout	    = 	200
	GGStatusLogoutAll	=	201
)