import QtQuick 2.1


Rectangle {
	id: screen
	width: 490; height: 320
	color: "black"
	Text {
		x: 8
		y: 10
		color: "white"
		font.pointSize: 20
		text: "Ввод:"
	}

	Text {
		x: 8
		y: 50
		color: "white"
		font.pointSize: 20
		text: "IP Адрес:"
	}

	TextInput {
		id: inputIP
		width: 240
		x : 120
		y : 10
		text: "127.0.0.1"
		font.pointSize: 20
		color: "blue"
		focus: true
		onAccepted: bridge.handleClick(inputIP, resultIP)
		Component.onCompleted: inputIP.selectAll()
	}

	Text {
		x: 120
		y: 50
		id: resultIP
		font.pointSize: 20
		color: "steelblue"
		text: "Введите IP адрес и нажмите ввод"
	}


	Text {
		anchors.horizontalCenter: parent.horizontalCenter
		anchors.bottom: parent.bottom
		text: "Выдает информацию о стандартном адресе"
		color: "steelblue"
		font.pointSize: 20
	}

}