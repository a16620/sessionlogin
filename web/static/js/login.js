function isValid(str) {
    if (str.length == 0) {
        return false;
    }

    return true;
}

function login(name, password) {
    return new Promise(function(resolve, reject) {
        $.ajax({
            url: "/login",
            data: {id: name, pw: password},
            type: "POST",
            enctype: "multipart/form-data",
            success: function(data) {
                if (data == "success") {
                    resolve(data);
                } else {
                    reject(data);
                }
            },
            error: function(request,status,error) {
                const msg = request.status+":"+request.responseText+":"+error;
                console.log(msg);
                reject(msg);
            }
        });
    });
}

