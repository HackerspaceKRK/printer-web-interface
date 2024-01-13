MD.Toolbar = function(){

  // tools left
  $("#tools_left .tool_button").on("click", function(){ 
    const mode = this.getAttribute("data-mode");

    if(mode == "qr") {
      const contents = prompt("Enter text to encode:", "");
      if(!contents) return;

      const size = 100;
      var svgNode = QRCode({
      
           msg :  contents
          ,dim :   600
          ,pad :   6
          ,mtx :   7
          ,ecl :  "H"
          ,ecb :   0
          ,pal : ["#000000", "#ffffff"]
          ,vrb :   1
      
      });
      const svg = new XMLSerializer().serializeToString(svgNode);
      svgCanvas.importSvgString(svg, true);
    }
    state.set("canvasMode", mode)
    if (mode === "shapelib") showShapeLib()
  });

  function setMode(mode) {
    $(".tool_button").removeClass("current");
    $("#tool_" + mode).addClass("current");
    $("#workarea").attr("class", mode);
    svgCanvas.setMode(mode);
  }

  function showShapeLib(){
    $("#tools_shapelib").show();
  }

  this.setMode = setMode;
}
