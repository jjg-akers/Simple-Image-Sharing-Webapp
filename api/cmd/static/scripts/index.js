window.onload=function(){
    // // modal open/close
    // $(".trigger").click(function() {
    //     $(this).toggleClass("modal")
    //     // $(".card-front").toggleClass("modal")
    //     // $(".card-front").toggleClass("card")
    //     // $(".card-from").toggleClass("show-modal")
    //     // $(".modal").toggleClass("show-modal", "slow");
    // });

    // $(".close-button").click(function() {
    //     $(".modal").toggleClass("show-modal", "slow");
    // });
    $('.hover-scale').hover(function() {
        $(this).addClass('transition');
    
    }, function() {
        $(this).removeClass('transition');
        // $(this).addClass('card-click');
    });

    $('.card-inner').click(function() {
        $(this).toggleClass('card-click');
        // $(this).addClass('transition');
    });
    // , function() {
    //     $(this).removeClass('transition');
    // });
}