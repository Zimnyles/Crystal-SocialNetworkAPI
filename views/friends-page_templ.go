// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "zimniyles/fibergo/views/layout"
import "zimniyles/fibergo/views/widgets"
import "zimniyles/fibergo/internal/models"

func FriendsPage(FriendPageCredentials models.FriendPageCredentials) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<main>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = layout.HeaderSmall().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = FriendsPageStyle().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<div class=\"page-container\"><div class=\"left-menu\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = widgets.LeftMenu().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</div><div class=\"main-content\"><div class=\"friends-header\"><div class=\"friendslist-header\"><!-- Заголовок списка друзей --><span>Ваши друзья</span></div><div class=\"friendsrequest-header\"><span>Запросы в друзья</span></div></div><div class=\"friends-content\"><div class=\"page-firendslist\" hx-swap=\"innerHTML\" hx-target=\"this\"><!-- Добавьте hx-target -->")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = widgets.FriendList(FriendPageCredentials.Friends).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, "</div><div class=\"page-friendsrequests\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = widgets.FriendRequestList(FriendPageCredentials.FriendRequests).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "</div></div></div></div></main>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return nil
		})
		templ_7745c5c3_Err = layout.Layout(layout.LayoutProps{
			Title:           "Друзья",
			MetaDescriptiom: "Ваши друзья",
		}).Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func FriendsPageStyle() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "<style>\r\n    .page-container {\r\n        display: flex;\r\n        width:1320px;\r\n        margin:0px auto;\r\n        box-sizing: border-box;\r\n    }\r\n\r\n   \r\n\r\n    .main-content {\r\n        margin-top: 15px;\r\n        flex-grow: 1;\r\n    }\r\n\r\n    .friends-header {\r\n        display: flex;\r\n        gap: 15px;\r\n        margin-bottom: 15px;\r\n    }\r\n\r\n    .friendslist-header {\r\n        \r\n        flex: 2;\r\n        height: 54px;\r\n        background: #222222;\r\n        border-radius: 20px;\r\n        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);\r\n        display: flex;\r\n        align-items: center;\r\n        text-align: center;\r\n        justify-content: center;\r\n        width: 100%;\r\n        color: var(--color-white);\r\n        text-decoration: none;\r\n    }\r\n\r\n    .friendsrequest-header {\r\n        min-width: 354px;\r\n        flex: 1;\r\n        height: 54px;\r\n        background: #222222;\r\n        border-radius: 20px;\r\n        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);\r\n        display: flex;\r\n        align-items: center;\r\n        text-align: center;\r\n        justify-content: center;\r\n        width: 100%;\r\n        color: var(--color-white);\r\n        text-decoration: none;\r\n    }\r\n\r\n    .friends-content {\r\n        display: flex;\r\n        gap: 15px;\r\n    }\r\n\r\n    .page-firendslist {\r\n        flex: 2;\r\n        \r\n    }\r\n\r\n    .page-friendsrequests {\r\n        flex: 1;\r\n        \r\n    }\r\n</style>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
