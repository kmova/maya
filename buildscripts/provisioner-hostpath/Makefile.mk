
# Specify the name for the binaries
PROVISIONER_HOSPATH=provisioner-hostpath

#Use this to build provisioner-hostpath
provisioner-hostpath:
	@echo "----------------------------"
	@echo "--> provisioner-hostpath    "
	@echo "----------------------------"
	@PNAME=${PROVISIONER_HOSPATH} CTLNAME=${PROVISIONER_HOSPATH} sh -c "'$(PWD)/buildscripts/build.sh'"

provisioner-hostpath-image: provisioner-hostpath
	@echo "-------------------------------"
	@echo "--> provisioner-hostpath image "
	@echo "-------------------------------"
	@cp bin/provisioner-hostpath/${PROVISIONER_HOSPATH} buildscripts/provisioner-hostpath/
	@cd buildscripts/provisioner-hostpath && sudo docker build -t openebs/provisioner-hostpath-ci:${IMAGE_TAG} --build-arg BUILD_DATE=${BUILD_DATE} . --no-cache
	@rm buildscripts/provisioner-hostpath/${PROVISIONER_HOSPATH}
