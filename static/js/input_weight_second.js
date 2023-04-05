window.onload = function (e) {
    liff.init(function (data) {
        initializeApp(data);
    });
};

function initializeApp(data) {

    // sendMessages call
    document.getElementById('submit').addEventListener('click', function () {
        var weight;

        weight = document.getElementById('weight').value;
        liff.sendMessages([{
            type: 'text',
            text: "答 " + weight
        }]).then(function () {
            window.alert("Message sent");
        }).catch(function (error) {
            window.alert("Error sending message: " + error);
        });
    });
}