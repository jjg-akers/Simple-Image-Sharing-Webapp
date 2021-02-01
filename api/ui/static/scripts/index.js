window.onload=function(){
    $('.hover-scale').hover(function() {
        $(this).addClass('transition');
    
    }, function() {
        $(this).removeClass('transition');
    });

    // display image
    $('.card-inner').click(function(e) {
        
        if ($(e.target).hasClass('trigger')){
            console.log("trigger clicked");
            $(e.target).closest('.card').next('.modal').toggleClass("show-modal");
            $('.card').toggle();
        } else{
            $(this).toggleClass('card-click');
        }
    });


    $("#search-form").submit(function(e) {
        tag = $("#tag-search-val").val().trim();

        console.log("tag: ", tag)

        if (tag === "") {
            $( "#invalid" ).text( "Invalid" ).show().fadeOut( 2000 );
            e.preventDefault();
            return
        }
    });

        
    //close modal
    $(".close-button").click(function(e) {
        $(this).parent().parent().toggleClass("show-modal");
        $('.card').toggle();
        
    });

    //save image link to clipboard
    $('.save-link').click(function(e){
        srcStr = $(this).parent().parent().parent().prev().find('img').attr('src');
        navigator.clipboard.writeText(srcStr);
        alert("Image Saved to Clipboard");
    });

}

//close modal
function windowOnClick(event){
    if(event.target === $(".modal")){
        $(".modal").toggleClass("show-modal");
    }
}

    // route handlers
// function callTagSearch(tag){
//         var query ="/search";
//         $.ajax({
//             method: 'GET',
//             url: query,
//             //dataType: 'json',
//             data: {
//                 "tag": tag
//                 // "reset": choice,
//             }
//         }).done(function(response) {

//             var json = JSON.stringify(eval("(" + response + ")"));

//             console.log("resp json: ", json)

//             var jsonObj = JSON.parse(json);
//         //console.log("call tag search success response ", response);
//             //i = JSON.eval(response);
//             //console.log("call tag search success response ", JSON.parse(response));
//             // var crypto = requery('crypto')
//             // var shasum = crypto.

//             if (jsonObj.length === "" ){
//                 displaySearchMissModal(tag);
//                 window.location.href = "http://localhost/";
//                 return;
//             }

//             handleImagesResponse(jsonObj);
//             return;
    
//         }).fail(function(xhr, status, error){
//             //console.log("error occured during callReset ajax: ", xhr.status);
            
//             console.log("there error is: ", error);
//             //console.log("status: ", status);

//             //console.log(JSON.parse(xhr.responseText));

//             return "";
//         });
    
//         setTimeout(function(){
//            return "";
//         }, 400);
// }

// function handleImagesResponse(images){
//         console.log("Handling images response");

//         // <div class="grid-item"></div>
//         //images-container
//         $("#images-container").replaceWith('<div class="grid-container" id="images-container"></div>');


//         images.forEach(element => {
//             str = `<div class="card">
//             <div class="card-inner">
//               <div class="card-front hover-scale">
//                     <img src=`;
//             str += element.url +">"
//             str += `</div>
//               <div class="card-back">
//                 <div class="back-contents">
//                   <p class="card-title">`;
//             str += element.Meta.Title + "</p>";
//             str += "<p>Tags: " + element.Meta.Tag + "</p>";
//             str += '<p><a href="#" class="trigger">This is a link</a></p>';
//             str += '<p>' + element.Meta.Description +'</p>';
//             str += `</div>
//               </div>
//             </div>
//           </div>`;

//             $("#images-container").append(str);
//         });
//     // $("#images-container").html(function(images){
//     //     var newHTML = ""
//     //     images.array.forEach(element => {
//     //         str = `<div class="card">
//     //         <div class="card-inner">
//     //           <div class="card-front hover-scale">
//     //                 <img src=`;
//     //         str += element.URI +">"
//     //         str += `</div>
//     //           <div class="card-back">
//     //             <div class="back-contents">
//     //               <p class="card-title">`;
//     //         str += element.Meta.Title + "</p>";
//     //         str += "<p>Tags: " + element.Meta.Tag + "</p>";
//     //         str += '<p><a href="#" class="trigger">This is a link</a></p>';
//     //         str += '<p>' + element.Meta.Description +'</p>';
//     //         str += `</div>
//     //           </div>
//     //         </div>
//     //       </div>`;
          

//     //     });
//     // });
// }

//     function displaySearchMissModal(tag){
//         //console.log("displaying reset modal")
//         //var player = gameState.turn.name;
    
//         //var str = "<span id=winner>" + player + " wants to Reset the game<br>";
//         var str = "<span id=modal_title> No images found for tag: " + tag + "RESET to confirm,<br>or CANCEL to return to game </span>";
//         str += '<div id="resetBtns">\
//                     <div align="center" class="trigger" id="confirm">\
//                         <span class="example_a">OK</span>\
//                     </div>\
//                 </div>';
        
//         $("#resultPane").html(str);
    
//         $("#resultModal").toggleClass("show-modal");

//         $(".close-button").click(function() {
//             // close modal and redirect
//             // window.location.href = "http://localhost/";
//             return;
//         });
//     }

