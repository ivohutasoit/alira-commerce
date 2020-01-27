function setActiveNavlink(page) {
    /*var current = sessionStorage.getItem("page");
    var link;
    if (current) {
        var currentPageLink = document.getElementById(current);
        if (currentPageLink) {
            currentPageLink.classList.remove("active")
        }
    }*/
    console.log(view)
    if (!view) view = "home";
    var link = document.getElementById(view);
    if (link) {
        link.classList.add("active");
        //sessionStorage.setItem("page", page);
    }
}