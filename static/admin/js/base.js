$(function () {

	$('.aside h4').click(function () {

		$(this).siblings('ul').slideToggle();
	})
})
$(function () {
	baseApp.init();
	$(window).resize(function(){
		baseApp.resizeIframe();
	})
})
var baseApp = {
	init: function () {
		this.initAside()
		this.confirmDelete()
		this.resizeIframe()
		this.changeStatus()
		this.changeNum()
	},
	initAside: function () {
		$('.aside h4').click(function () {
			$(this).siblings('u1').slideToggle();
		})
	},
	// 设置右侧iframe高度
	resizeIframe: function () {
		$('#rightMain').height($(window).height() - 80)
	},
	// 删除提示
	confirmDelete: function () {
		$('.delete').click(function () {
			var flag = confirm('确定要删除吗？')
			return flag
		})
	},
	// 改变状态
	changeStatus:function(){
		$(".chStatus").click(function(){
			var id=$(this).attr("data-id")
			var table=$(this).attr("data-table")
			var field=$(this).attr("data-field")
			var el=$(this)
			$.get("/admin/changeStatus",{id:id,table:table,field:field},function(response){
				if(response.success){
					if(el.attr("src").indexOf("yes")!=-1){
						el.attr("src","/static/admin/images/no.gif")
					}else{
						el.attr("src","/static/admin/images/yes.gif")
					}

				}
			})

		})
	},
	// 改变排序值
	changeNum:function(){
		$(".chSpanNum").click(function(){
			//1、获取el里面的值
			var id=$(this).attr("data-id")
			var table=$(this).attr("data-table")
			var field=$(this).attr("data-field")
			var num=$(this).html().trim()
			var SpanEl=$(this)
			//2、创建一个文本框
			var input=$("<input style='width:60px' value='' />")
			//3、把input放在SpanEl里
			$(this).html(input)
			//4、让input获取焦点，为input赋值 这里focus是焦点
			$(input).trigger("focus")
			$(input).val(num)
			//5、点击input时阻止冒泡
			$(input).click(function(e){
				e.stopPropagation()
			})
			//6、鼠标离开时给span赋值，并且触发一个ajax请求
			$(input).blur(function(){
				var inputNum=$(this).val()
				SpanEl.html(inputNum)
				// 触发ajax请求，更新数据库
				$.get("/admin/changeNum",{id:id,table:table,field:field,num:inputNum},function(response){
					// if(response.success){
					// 	if(el.attr("src").indexOf("yes")!=-1){
					// 		el.attr("src","/static/admin/images/no.gif")
					// 	}else{
					// 		el.attr("src","/static/admin/images/yes.gif")
					// 	}
					// }
					console.log(response)
				})
			})
		})
	}
}