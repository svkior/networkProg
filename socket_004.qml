import QtQuick 2.1


Rectangle {
	id: screen
	width: 800; height: 600
	color: "black"
	Text {
		x: 8
		y: 10
		color: "white"
		font.pointSize: 20
		text: "Ввод:"
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
		onAccepted: bridge.handleClick(inputIP)
		Component.onCompleted: inputIP.selectAll()
	}


	ListView{
		id: listView1
		anchors.right: parent.right
		anchors.left: parent.left
		anchors.leftMargin: 10
		anchors.bottom: bottomBar.top
		anchors.bottomMargin: 30
		anchors.top: inputIP.bottom
		anchors.topMargin: 30
		//anchors.top: parent.top 
		model: logs.len
		delegate: Row {
			id: row1
			height: 30
			Rectangle { 
				id: boxrect
				width: 40; height: 20; color: "gold"; radius: 20
				Text {
					anchors.centerIn: parent
					text: index
				}
			}
			Rectangle { width: 20; height:40; color: "black"}
			Text {
				text: logs.record(index)
				color: "gold"
				anchors.verticalCenter: boxrect.verticalCenter
				font.bold: true
			}
		}
	}


	Text {
		id: bottomBar
		anchors.horizontalCenter: parent.horizontalCenter
		anchors.bottom: parent.bottom
		text: "Выводит все IP адреса для данного домена"
		color: "steelblue"
		font.pointSize: 20
	}

}