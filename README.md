# gofaces
go语言的边缘计算摄像头设备

环境搭建(开发与测试环境)：

1、构建dlib19.17并安装

    1.1、安装依赖:

        $sudo apt-get install libdlib-dev libopenblas-dev libjpeg-turbo8-dev
        $sudo apt-get install libboost-all-dev

    1.2、构建并安装：
        $cd dlib19.17
        $mkdir build
        $cd build
        $cmake ..
        $cmake --build . --config Release
        $cd dlib
        $make install
2、修改dlib的pkg-config文件
    $vi /usr/local/lib/pkgconfig/dlib-1.pc
    
    libdir=/usr/local/lib
    includedir=/usr/local/include
    
    Name: dlib
    Description: Numerical and networking C++ library
    Version: 19.17.0
    Libs: -L${libdir} -ldlib -lblas -llapack
    Cflags: -I${includedir}
    Requires: libpng
3、ffmpeg 安装
  
    #sudo apt-get install ffmpeg
    
环境搭建（运行环境）：
    
4、接口调用
    
    IP = 192.168.0.162
    PORT = 8080
    FACENAME = "要建模的人脸的名字"
    构建人脸模型
    GET：http://$(IP):$(PORT)/api/v1/buildFaceModle?facename=$(FACENAME)
    识别人脸：
    GET：http://$(IP):$(PORT)/api/v1/classifyFace

环境搭建（运行环境）

1、构建交叉编译环境
    
    1.1、安装arm-linux-gnueabihf-交叉编译工具链：
   https://developer.arm.com/tools-and-software/open-source-software/developer-tools/gnu-toolchain/gnu-a/downloads
   
   下载AArch64 GNU/Linux big-endian target (aarch64_be-linux-gnu)工具包并安装到宿主机中
 
2、在目标机器安装环境

    在目标环境中安装必要的库
    
    将书梅派中的/lib /usr文件夹整个拷贝到rootfs根目录下
    rsync -rl --delete-after --safe-links root@192.168.184.199:/{lib,usr} $HOME/arm-dlib/rootfs
  
3、编译dlib
    
    下载并解压dlib
    创建CMAKE_TOOLCHAIN_FILE：
    在dlib目录下新建cmake脚本pi.cmake，内容如下
    #设置系统属性
    SET(CMAKE_SYSTEM_NAME Linux)
    SET(CMAKE_SYSTEM_VERSION 1)
    设置c/cxx编译器路径
    SET(CMAKE_C_COMPILER $ENV{HOME}/arm-dlib/OrangePiH3_toolchain/bin/arm-linux-gnueabihf-gcc)
    SET(CMAKE_CXX_COMPILER $ENV{HOME}/arm-dlib/OrangePiH3_toolchain/bin/arm-linux-gnueabihf-g++)
    #设置根目录查找范围
    SET(CMAKE_FIND_ROOT_PATH $ENV{HOME}/arm-dlib/rootfs)
    SET(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
    SET(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
    SET(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
    编译dlib
    mkdir build
    cd build
    cmake -DCMAKE_TOOLCHAIN_FILE=/home/eagle/arm-dlib/dlib-19.17/OrangePI.cmake ../
    cmake --build . --config Release
    
4、安装ffmpeg
    
    
