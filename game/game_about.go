package game

//  order type key          value   description
const (    
    OT_Msg              =   10000   // message order
    OT_MsgPerson        =   11000   // msg from client to client
    OT_MsgPersonGroup   =   11001   // msg from client to group created by the client
    OT_MsgNear          =   12000   // msg from client to group that contains near clients
    OT_MsgGroup         =   12001   // msg from client to a group create by the system
    OT_MsgAll           =   12002   // msg from client to all
    OT_MsgSystem        =   13000   // normal msg from system to client to clients
    OT_MsgSystemInner   =   13001   // inner msg from system to client to clients
    OT_MsgSystemPerson  =   13002   // feedback msg from system to client

    OT_Skill            =   20000   // skill order
    OT_SkillSingle      =   21000   // single target skill with damage
    OT_SkillSingleK     =   22000   // single target skill with control
    OT_SkillNear        =   25000   // near targets skill with damage
    OT_SkillNearK       =   27000   // near targets skill with control
    
    OT_Action           =   30000   // action order
    OT_ActionDrug       =   31000   // take some drug
    OT_ActionMove       =   32000   // client actions that relative to move 

    OT_Data             =   40000   // game data signal order
    OT_DataPerson       =   40001   // update personal game data
    OT_DataAll          =   40002   // update clients game data
    OT_DataPersonLogin  =   40003   // feedback to the other clients when a client login
)  

//  id type key             value   description
const (
    ITPerson            =   10000   // client id group
    ITGroup             =   20000   // clients id group
    ITSystem            =   30000   // system id group

    IDSYSTEM    int64   =   30000   // system id
    IDPERSON    int64   =   10000   // client id
)

const (
    GStatusLogin 	    =	10000
	GStatusInGame	    =   10001
    GStatusLogout	    = 	20000
    GStatusErrorLogout  =   20001
	GGStatusLogoutAll	=	20002
)