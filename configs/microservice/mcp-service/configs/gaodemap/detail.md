功能描述
file_list	基础文件信息-获取文件列表	获取用户网盘中指定目录下的文件列表。
file_doc_list	基础文件信息-获取文档列表	获取用户指定目录下的文档列表。
file_image_list	基础文件信息-获取图片列表	获取用户指定目录下的图片列表。
file_video_list	基础文件信息-获取视频列表	获取用户指定目录下的视频列表。
make_dir	文件管理-创建文件夹	创建文件夹。
file_copy	文件管理-复制	对指定的文件进行复制操作。
file_del	文件管理-删除	对指定的文件进行删除操作。
file_move	文件管理-移动	对指定的文件进行移动操作。
file_rename	文件管理-重命名	对指定的文件进行重命名操作。
file_upload_stdio	文件上传-上传本地文件	将用户本地文件上传存储在网盘的云端文件。
因需对本地文件读取，上传本地文件仅支持stdio模式。
file_upload	文件上传-通过URL或文本上传	以URL或text文本的方式上传文件到网盘，返回保存为的网盘文件的路径
file_download	文件下载	将用户存储在网盘的云端文件下载到本地。
file_search	文件搜索	获取用户指定目录下，包含指定关键字的文件列表。
user_info	用户基础信息-获取已鉴权用户信息	获取用户的基本信息。
get_quota	用户基础信息-获取网盘容量信息	获取用户的网盘空间的使用情况。
服务密钥（AK）获取流程

一、企业开发者接入
需要完成基本流程才能使用开放平台的各种能力以及服务，通用的操作流程如下图所示：

1. 注册登录百度账号
在使用服务之前，开发者需要在百度网盘开放平台注册百度账号，具体流程请参考【注册登录百度账号】https://pan.baidu.com/union/doc/Ll0ffxg47

2. 实名认证
账号注册完后，需要完成实名认证才能享受网盘的各类能力和服务，具体流程请参考【实名认证】https://pan.baidu.com/union/doc/ll0g1s7jb

3. 创建应用
完成实名认证后，开发者需要先在开放平台的控制台中完成应用的创建以获得关键接入信息，具体流程请参考【创建应用】。

4. 获取服务密钥（AK）
目前开放平台的服务均基于百度OAuth体系构建，如果需要接入开放平台的服务，需要先为您的应用接入授权，获取Access Token，具体流程请参考【接入授权】。



二、个人用户（限时体验）

1. 注册登录百度账号
在使用服务之前，开发者需要注册百度账号，具体流程请参考【注册登录百度账号】https://pan.baidu.com/union/doc/Ll0ffxg47

2. 获取服务密钥（AK）
目前开放平台的服务均基于百度OAuth体系构建，如果需要接入开放平台的服务，需要先为您的账号接入授权，获取Access Token。具体流程如下：

2.1 发起授权请求
【前往该链接发起授权请求】
https://openapi.baidu.com/oauth/2.0/authorize?response_type=token&client_id=QHOuRXiepJBMjtk0esLhrPoNlQyYd0mF&redirect_uri=oob&scope=basic,netdisk

2.2 获取Access Token
点击上述链接后，跳转至授权页面，用户点击“授权”按钮后跳转至Access Token链接。复制链接中Access Token部分即可。


3. 授权管理
如您需要进行应用解绑，可在 【授权管理】 管理网盘MCP Server的授权，对授权过的应用做权限“设置”、“解除关联”等操作。
https://passport.baidu.com/v2/?login&u=https%3A%2F%2Fpassport.baidu.com%2Fv6%2FappAuthority

使用百度网盘MCP Server
获取到Access Token 凭证就可以调用网盘MCP Server，访问用户信息以及授权资源了


注意：百度网盘开放平台目前仅限企业开发者注册调用，本模块提供的密钥信息仅供测试体验，将不定期变更密钥，请勿用于正式环境。
