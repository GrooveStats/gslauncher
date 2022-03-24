Unicode true
!include LogicLib.nsh
!include x64.nsh

!define APPNAME "GrooveStats Launcher"

!define VERSIONMAJOR 1
!define VERSIONMINOR 3
!define VERSIONBUILD 0

RequestExecutionLevel admin

InstallDir "$PROGRAMFILES\${APPNAME}"

Name "${APPNAME}"
Icon logo.ico
outFile ../dist/gslauncher-windows-setup.exe


page components
page directory
Page instfiles

!macro VerifyUserIsAdmin
UserInfo::GetAccountType
Pop $0
${If} $0 != "admin"
        MessageBox mb_iconstop "Administrator rights required!"
        SetErrorLevel 740 ;ERROR_ELEVATION_REQUIRED
        Quit
${EndIf}
!macroend

function .onInit
	SetShellVarContext all
	!insertmacro VerifyUserIsAdmin
functionEnd

section "GrooveStats Launcher"
	SectionIn RO
	SetOutPath $INSTDIR

	${If} ${RunningX64}
		File /oname=gslauncher.exe ../dist/gslauncher-windows-amd64.exe
	${Else}
		File /oname=gslauncher.exe ../dist/gslauncher-windows-i386.exe
	${EndIf}
	File logo.ico

	writeUninstaller "$INSTDIR\uninstall.exe"

	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayName" "${APPNAME}"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "UninstallString" "$\"$INSTDIR\uninstall.exe$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "QuietUninstallString" "$\"$INSTDIR\uninstall.exe$\" /S"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "InstallLocation" "$\"$INSTDIR$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayIcon" "$\"$INSTDIR\logo.ico$\""
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "Publisher" "GrooveStats"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "URLUpdateInfo" "https://github.com/GrooveStats/gslauncher/releases"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "URLInfoAbout" "https://github.com/GrooveStats/gslauncher"
	WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayVersion" "${VERSIONMAJOR}.${VERSIONMINOR}.${VERSIONBUILD}"
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "VersionMajor" ${VERSIONMAJOR}
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "VersionMinor" ${VERSIONMINOR}
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "NoModify" 1
	WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "NoRepair" 1
sectionEnd

section "Start Menu Shortcut"
	CreateShortCut "$SMPROGRAMS\${APPNAME}.lnk" "$INSTDIR\gslauncher.exe" "" "$INSTDIR\logo.ico"
sectionEnd

section "Desktop Shortcut"
	CreateShortCut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\gslauncher.exe" "" "$INSTDIR\logo.ico"
sectionEnd

# Uninstaller

function un.onInit
	SetShellVarContext all

	MessageBox MB_OKCANCEL "Uninstall ${APPNAME}?" IDOK next
	Abort

	next:
	!insertmacro VerifyUserIsAdmin
functionEnd

section "un.uninstall"
	delete "$SMPROGRAMS\${APPNAME}.lnk"
	delete "$DESKTOP\${APPNAME}.lnk"

	delete $INSTDIR\gslauncher.exe
	delete $INSTDIR\logo.ico
	delete $INSTDIR\uninstall.exe
	rmDir $INSTDIR

	DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}"
sectionEnd
