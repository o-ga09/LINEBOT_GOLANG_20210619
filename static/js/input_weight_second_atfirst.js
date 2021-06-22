window.onload = function (e) {
    liff.init(function (data) {
        initializeApp(data);
    });
};

function initializeApp(data) {

    // sendMessages call
    document.getElementById('submit').addEventListener('click', function () {
        var weight;
        var height;

        weight = document.getElementById('weight').value;
        height = document.getElementById('height').value;
        liff.sendMessages([{
            type: 'text',
            text: "答 " + weight + "," + height
        }]).then(function () {
            window.alert("Message sent");
        }).catch(function (error) {
            window.alert("Error sending message: " + error);
        });
    });
}