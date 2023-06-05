## 定制主题
您可以自定义文件浏览器安装，方法是将其名称更改为所需的任何其他名称，添加全局自定义样式表，并根据需要使用自己的标识。为此，可以更改三个配置选项：
+ 名称：这是将显示在登录和注册页面上的实例名称。这不会替换边栏中的版本消息。
+ 禁用外部链接：这将禁用任何外部链接（本文档的外部链接除外）。
+ 文件夹：是可以包含两个项目的目录的路径：

自定义.css，包含要应用于安装的样式。
img 一个目录，其文件可以替换 在应用程序中。
可以使用以下命令通过 CLI 界面设置这些选项：

```shell

filebrowser config set --branding.name "My Name" \

    --branding.files "/abs/path/to/my/dir" \
    
    --branding.disableExternal
```

或者可以在“设置”→“全局设置”中的“品牌目录路径”下进行设置。
如果使用Docker，那么请记住绑定此目录，例如`'/home/username/containers/filebrowser/branding：/branding'`,要识别自定义图标，您需要创建 img 和  img/icons  目录，并将 svg 放在 branding/img 目录中： ​
要替换图标，您需要将其放在 img/icons 目录中，但也请注意还需要其他一些 PNG 图标类型（请参阅上面的默认标识链接），因为浏览器通常会使用可用的最高分辨率选项（至少 16x16 和 32x32 选项）。您可以使用 以从基本映像生成这些映像。
## 命令运行器
命令运行器是一项功能，使您能够在特定事件之前或之后执行您想要的任何 shell 命令。现在，这些是操作：

+ Copy
+ Rename
+ Upload
+ Delete
+ Save

此外，在执行为这些挂钩设置的命令期间，将有一些环境变量可用于帮助您执行命令：

具有更改文件的完整绝对路径的 FILE。

+ SCOPE 与用户范围的路径。
+ 带有事件名称的 TRIGGER。
+ USERNAME 为用户的用户名。
+ DESTINATION 到目的地的绝对路径。仅用于复制和重命名。

此时，您可以通过命令行界面编辑命令，使用以下命令（请检查标志 --help 以了解更多信息）：

+ 文件浏览器命令添加 before_copy "echo $FILE"
+ 文件浏览器 cmds rm before_copy 0
+ 文件浏览器命令 ls

或者您可以使用 Web 界面通过设置 → 全局设置来管理它们。
