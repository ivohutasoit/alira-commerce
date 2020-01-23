function setActiveNavlink(page) {
    /*var current = sessionStorage.getItem("page");
    var link;
    if (current) {
        var currentPageLink = document.getElementById(current);
        if (currentPageLink) {
            currentPageLink.classList.remove("active")
        }
    }*/
    if (!page) page = "home";
    var link = document.getElementById(page);
    if (link) {
        link.classList.add("active");
        //sessionStorage.setItem("page", page);
    }
}