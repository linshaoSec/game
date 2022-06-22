package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type obj struct {
	address int
}

var (
	flag1 int
	iSsuoxie int
	iSfeitian int
	pid int64
	handle uintptr
	modelName="ResonanceEAE-Win64-Shipping.exe"
	kernel32=syscall.MustLoadDLL("kernel32.dll")
	OpenProcess=kernel32.MustFindProc("OpenProcess")
	ReadProcessMemory=kernel32.MustFindProc("ReadProcessMemory")
	WriteProcessMemory=kernel32.MustFindProc("WriteProcessMemory") //写入内存
	Psapi=syscall.MustLoadDLL("Psapi.dll")
	EnumProcessModules =Psapi.MustFindProc("EnumProcessModules")
	GetModuleBaseNameA=Psapi.MustFindProc("GetModuleBaseNameA")


	xueAddr int64
	hightAddr int64
	hightValue int=1170409677//最佳飞天高度
	xueValue int=1272549262//最佳血量
	currenxueValue int
	currenthightValue int
	tmphight int
	//xueRVA1=[...]int64{0x04CF3798,0x68,0x110,0xE0,0x260,0x1A0,0x90,0xB8}
	xueRVA1=[...]int64{0x04CF3798,0x120,0x110,0xF8,0x20,0x1A0,0x90,0xB8}
	//血量其他偏移路径
	// 0x04CF3798,0x120,0x110,0xF8,0x20,0x1A0,0x90,0xB8
	// 04DB5400 400 8 78 A0 1A0 90 B8
	hightRVA1=[...]int64{0x04DB4DE8,0x50,0x20,0x220,0x250,0x1A0,0,0x1D8}
	//高度其他偏移路径
	//048B0710,10,C8,130,260,140,C0,1D8
	//04CF3798,308,110,E0,250,1A0,0,1D8
	//04DB4DE8,50,20,220,250,1A0,0,1D8
)

func main(){

	//打开获取句柄
	handle=getHandle(modelName)
	//获取模块基址
	BaseAddress:=GetProcessMoudleBase(handle,modelName)	//fmt.Printf("[+] 基址%X\n",BaseAddress)
	//fmt.Printf("==========================================================\n\t\t==linshao--功能面板==\t\t\n")
	fmt.Printf("[*] 血量：使用的RVA偏移量链 -》%X\n",xueRVA1)
	fmt.Printf("[*] 高度：使用的RVA偏移量链 -》%X\n",hightRVA1)
	xueAddr,currenxueValue=readXue(handle,BaseAddress,xueRVA1)
	hightAddr,currenthightValue=readHight(handle,BaseAddress,hightRVA1)

	for true{

		fmt.Printf("==>状态\n")
		fmt.Printf("[!]当前血量--->[%d](%f%%)\n",currenxueValue,(float32(currenxueValue)/1120403456)*100)
		fmt.Printf("[!]当前高度--->[%d]\n",currenthightValue)
		if iSsuoxie==0{
			fmt.Printf("锁血[关闭]\n")
		}
		if iSsuoxie==1{
			fmt.Printf("锁血[开启]\n")
		}
		if iSfeitian==0{
			fmt.Printf("飞天[关闭]\n")
		}
		if iSfeitian==1{
			fmt.Printf("飞天[开启]\n")
		}

		fmt.Printf("\n-----------------------------------------------" +
						  "\n[*] linshao：请输入功能号\n[*] 1.锁血 2.飞天 0退出->")
		var gn int=-1
		fmt.Scanln(&gn)
		if gn==1{
			fmt.Printf("[*] 1开启/0关闭->")
			var choice int
			fmt.Scanln(&choice)
			if choice==1{
				iSsuoxie=1
				currenxueValue=xueValue
				go suoxie(handle,xueAddr,&currenxueValue,&iSsuoxie)
			}
			if choice==0 {
				iSsuoxie=0
			}

		}
		if gn==2{
			fmt.Printf("[*] 1推荐飞天/2缓步升高/3缓步降落/0直接落地->")


			var choice int
			fmt.Scanln(&choice)

			if choice==1{
				currenthightValue=hightValue//使用推荐高度
				if iSfeitian==0{
					iSfeitian=1
					go feiTian(handle,hightAddr,&currenthightValue,&iSfeitian)
				}

			}
			if choice==2 {
				//如果没有启动就先飞起来
				if iSfeitian==0 {
					iSfeitian=1
					go feiTian(handle,hightAddr,&currenthightValue,&iSfeitian)
				}

					for true{

						fmt.Printf("[回车]升高 [0]:退出升高->")
						var tmpplus int=-1
						fmt.Scanln(&tmpplus)
						if tmpplus==-1{
							currenthightValue=currenthightValue+30000//高度++
						}
						if tmpplus==0{
							break
						}
					}
			}
			if choice==3 {
				//如果没有启动就先飞起来
				if iSfeitian==0 {
					iSfeitian=1
					go feiTian(handle,hightAddr,&currenthightValue,&iSfeitian)
				}

				for true{
					fmt.Printf("[回车]下降 [0]:退出下降->")
					var tmpplus int=-1
					fmt.Scanln(&tmpplus)
					if tmpplus==-1{
						currenthightValue=currenthightValue-40000//高度++
					}
					if tmpplus==0{
						break
					}
				}
			}
			if choice==0 {
				iSfeitian=0
				fmt.Printf("[+]飞天后自动开启锁血，免得摔死\n")
				if iSsuoxie==0{
					iSsuoxie=1
					go suoxie(handle,xueAddr,&xueValue,&iSsuoxie)
				}
			}
		}
		if gn==0{
			syscall.Exit(0)
		}

		fmt.Printf("==========================================================\n")
	}
}
func getHandle(modelname string)uintptr{
	pid=getProcPid(modelName)
	fmt.Println(pid)
	return getModelHandle(pid)
}


//读取血量
func readXue(handle uintptr,BaseAddress int64,RVA2 [8]int64) (int64,int){
	return getAddvalueByBassaddAndRVA(handle ,BaseAddress ,RVA2 )
}
//读取高度
func readHight(handle uintptr,BaseAddress int64,RVA2 [8]int64) (int64,int){
	return getAddvalueByBassaddAndRVA(handle ,BaseAddress ,RVA2 )
}
//锁血
func suoxie(handle uintptr,xueAddr int64,mxueValue *int,miSsuoxie *int){
	setValueByAddvalue(handle,xueAddr,mxueValue,miSsuoxie)
}
//飞天
func feiTian(handle uintptr,hightAddr int64,mhightValue *int, miSfeitian *int){
	setValueByAddvalue(handle,hightAddr,mhightValue,miSfeitian)
}

//循环修改某个地址的值
func setValueByAddvalue(handle uintptr,theAddr int64,theValue *int,flag *int){ //必须传指针,不能取值，否侧携程无法确定对应的变量具体值

	defer func() {
		err:=recover()
		if err!=nil{
			fmt.Println("成功捕获一个异常",err)
		}
	}()
	for *flag!=0{
		value:=0
		ReadProcessMemory.Call(handle, uintptr(theAddr),uintptr(unsafe.Pointer(&value)),unsafe.Sizeof(value),0)
		if value!=*theValue{
			WriteProcessMemory.Call(handle,uintptr(theAddr),uintptr(unsafe.Pointer(theValue)),unsafe.Sizeof(theValue))
		}
		time.Sleep(10*time.Millisecond)
	}
}

//通过基址+RVA链读取到指针地址的值
func getAddvalueByBassaddAndRVA(handle uintptr,BaseAddress int64,RVA2 [8]int64) (int64,int){
	var CurrentAddress int64
	//fmt.Println("模块名",modelName)
	CurrentAddress=BaseAddress
	var value int
	var tmp64 int64
	for _,x:=range RVA2{
		value=0
		tmp64=CurrentAddress+x
		ReadProcessMemory.Call(handle, uintptr(tmp64),uintptr(unsafe.Pointer(&value)),unsafe.Sizeof(value),0)
		//fmt.Printf("-->访问0x%X >> (RVA) %X -> (0x%X) %X\n",CurrentAddress,x,tmp64,value)
		CurrentAddress=int64(value) //把偏移后获取到的实际地址值放入下一次循环计算的基地址
	}
	return tmp64,value
}

//获取进程基址
func GetProcessMoudleBase( hProcess uintptr, moduleName string)int64{
	//异常处理
	defer func() {
		//捕获异常
		err := recover()
		if err != nil {
			fmt.Println("[!]发生异常")
		}
	}()
	// 遍历进程模块,
	var hModel =[10000]int64{0}
	var lpcbNeeded int=0 //将所有模块句柄存储在 lphModule 数组中所需的字节数
	var cb= int(unsafe.Sizeof(hModel))  //lphModule 数组的大小，以字节为单位。
	isok,_,_:=EnumProcessModules .Call(hProcess, uintptr(unsafe.Pointer(&hModel)), uintptr(cb), uintptr(unsafe.Pointer(&lpcbNeeded)))
	num:=lpcbNeeded/int(unsafe.Sizeof(hModel[0]))
	//fmt.Println("----------lpcbNeeded所需字节",lpcbNeeded,"当前",cb,"需要长度",num)
	if isok<=0 {
		fmt.Println("[!] 枚举模块失败")
	}
	fmt.Printf("[+] 枚举进程模块成功,共%d个\n",num)//,hModel)
	tmp:=[50]byte{}
	a:=""
	for i:=0;i<num;i++{
		GetModuleBaseNameA.Call(hProcess, uintptr(hModel[i]),uintptr(unsafe.Pointer(&tmp)),50)
		for _,v:=range tmp{
			if(v==0){continue}else {a+=string(v)}
		}
		if(strings.EqualFold(moduleName, a)) {
			fmt.Printf("[+] find! 模块名字%s  地址：0x%X\n", a, hModel[i])
			return int64(hModel[i])
		}
		fmt.Printf(" > %s  \t--->: 0x%X",a, hModel[i])
		fmt.Printf(" \n ")
		a=""
		tmp=[50]byte{}
	}
	return 0
}

//调用wmic获取程序的pid
func getProcPid(PROCESS string)int64{
	task:=exec.Command("cmd","/c","wmic", "process", "get", "name,","ProcessId","|","findstr",PROCESS)
	data2, _ := task.CombinedOutput()
	res:=strings.Split(string(data2),"\n")[0]//取第一行程序结果
	re := regexp.MustCompile("[0-9]+")
	pid,_:=strconv.ParseInt(re.FindAllString(res,-1)[1],10,64)
	return pid
}

//根据pid打开模块返回句柄
func getModelHandle(_pid int64)uintptr{
	hand,_, err :=OpenProcess.Call(2097151, uintptr(0), uintptr(_pid))
	fmt.Println("[+] 打开目标进程成功",hand)
	checkErr(err)
	if hand<=0{fmt.Printf(" 打开句柄失败",hand)
		syscall.Exit(1)}
	return hand
}

//检查错误
func checkErr(err error){
	if err!=nil {
		if err.Error() != "The operation completed successfully." {
			fmt.Println("报错：",err.Error())
			syscall.Exit(1)
			return
		}
	}
}
