# gofaces
go语言的边缘计算摄像头设备

环境搭建：

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
    
4、接口调用
    
    构建人脸模型
    GET：http://$(IP):$(PORT)/api/v1/buildFaceModle?facename=$(FACENAME)
    识别人脸：
    GET：http://$(IP):$(PORT)/api/v1/classifyFace

    
    
    
    
