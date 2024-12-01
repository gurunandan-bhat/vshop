(function () {
    let form = document.getElementById("add-to-cart");
    if (form === null) {
        return;
    }

    form.addEventListener("submit", (evt) => {
        evt.preventDefault();

        let idElem = document.getElementById("product-id");
        let iProdID = null;
        if (idElem) {
            iProdID = idElem.value;
        }

        let qtyElem = document.getElementById("quantity");
        let iQty = null;
        if (qtyElem) {
            iQty = qtyElem.value;
        }

        fetch("http://localhost:3000/add-to-cart", {
            method: "POST",
            body: JSON.stringify({
                iProdID: iProdID,
                iQty: iQty,
            }),
            headers: {
                "Content-type": "application/json charset-UTF-8",
            },
        })
            .then((response) => {
                if (response.ok) {
                    return response.text;
                }
                throw new Error("Something went wrong");
            })
            .then((data) => {
                console.log("Success: ", data);
            })
            .catch((err) => {
                console.log(err);
            });
    });
})();
