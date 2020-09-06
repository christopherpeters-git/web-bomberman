let index = 0;
let fruits = ["/media/img1.JPG", "/media/img2.JPG", "/media/img3.JPG"];
let slideshow = document.querySelector(".slideShow");

function slideRight() {
    index++;
    if (index >= fruits.length) {
        index = 0;
    }
    $(".slideShow")
        .stop()
        .animate({ opacity: 0 }, 500, function () {
            $(this)
                .css({ "background-image": "url('" + fruits[index] + "')" })
                .animate({ opacity: 1 }, { duration: 500 });
        });
    console.log(index);
}

function slideLeft() {
    index--;
    if (index < 0) {
        index = fruits.length - 1;
    }
    $(".slideShow")
        .stop()
        .animate({ opacity: 0 }, 500, function () {
            $(this)
                .css({ "background-image": "url('" + fruits[index] + "')" })
                .animate({ opacity: 1 }, { duration: 500 });
        });
    console.log(index);
}
