window.onload = function (e) {
    liff.init(function (data) {
        initializeApp(data);
    });
};

function initializeApp(data) {

    // sendMessages call
    document.getElementById('submit').addEventListener('click', function () {
        var label;
        var category;
        var total;

        label = document.getElementById('label').value;
        category = document.getElementById('category').value;
        total = document.getElementById('total').value;

        weight = document.getElementById('weight').value;
        liff.sendMessages([{
            type: 'text',
            text: "収支 " + label + "," + category + "," + total
        }]).then(function () {
            window.alert("Message sent");
        }).catch(function (error) {
            window.alert("Error sending message: " + error);
        });
    });
}