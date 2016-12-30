/**
 * Created by vitold on 30.12.16.
 */
"use strict";

$(function () {
	$("#content").bind("input change", function () {
		$.post("/gethtml", {md: $("#content").val()},
			function (response) {
			 $("#md_html").html(response.html)
		});
	});
})