import QtQuick 2.1


Rectangle {
	id: scrollbar
	width: 15
	radius: 3
	color: "#333333"
	border.color: "#000000"
	property real position: handle1.y / (height - handle1.height)
	Rectangle{
		id: handle1
		width: scrollbar.width
		height: 50
		color: "#d9d8d8"
		MouseArea{
			anchors.fill: parent
			drag.target: handle1

			drag.axis: Drag.YAxis
			drag.minimumY: 0
			drag.maximumX: scrollbar.height - handle1.height
		}
	}
}