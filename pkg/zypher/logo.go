package zypher

import "fmt"

const logo = `
███████╗██╗   ██╗██████╗ ██╗  ██╗███████╗██████╗ 
╚══███╔╝╚██╗ ██╔╝██╔══██╗██║  ██║██╔════╝██╔══██╗
  ███╔╝  ╚████╔╝ ██████╔╝███████║█████╗  ██████╔╝
 ███╔╝    ╚██╔╝  ██╔═══╝ ██╔══██║██╔══╝  ██╔══██╗
███████╗   ██║   ██║     ██║  ██║███████╗██║  ██║
╚══════╝   ╚═╝   ╚═╝     ╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝
`

func PrintLogo() {
	fmt.Print(logo)
}
