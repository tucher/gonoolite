package gonoolite

type EnumCMD int

const (
	Off              EnumCMD = 0
	Bright_Down      EnumCMD = 1
	On               EnumCMD = 2
	Bright_Up        EnumCMD = 3
	Switch           EnumCMD = 4
	Bright_Back      EnumCMD = 5
	Set_Brightness   EnumCMD = 6
	Load_Preset      EnumCMD = 7
	Save_Preset      EnumCMD = 8
	Unbind           EnumCMD = 9
	Stop_Reg         EnumCMD = 10
	Bright_Step_Down EnumCMD = 11
	Bright_Step_Up   EnumCMD = 12
	Bright_Reg       EnumCMD = 13
	Bind             EnumCMD = 15
	Roll_Colour      EnumCMD = 16
	Switch_Colour    EnumCMD = 17
	Switch_Mode      EnumCMD = 18
	Speed_Mode_Back  EnumCMD = 19
	Battery_Low      EnumCMD = 20
	Sens_Temp_Humi   EnumCMD = 21
	Temporary_On     EnumCMD = 25
	Modes            EnumCMD = 26
	Read_State       EnumCMD = 128
	Write_State      EnumCMD = 129
	Send_State       EnumCMD = 130
	Service          EnumCMD = 131
	Clear_memory     EnumCMD = 132
)

type EnumMode int

const (
	TX       EnumMode = 0
	RX       EnumMode = 1
	FTX      EnumMode = 2
	FRX      EnumMode = 3
	SVC      EnumMode = 4
	FWUpdate EnumMode = 5
)
