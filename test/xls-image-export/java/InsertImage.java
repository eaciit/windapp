package com.eaciit.wfdemo;
 
import java.io.FileInputStream;
import java.io.FileOutputStream;
import java.io.InputStream;
 
import org.apache.poi.xssf.usermodel.XSSFWorkbook;
import org.apache.poi.ss.usermodel.ClientAnchor;
import org.apache.poi.ss.usermodel.CreationHelper;
import org.apache.poi.ss.usermodel.Drawing;
import org.apache.poi.ss.usermodel.Picture;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.util.IOUtils;
 
 
public class InsertImage {
 
 public static void main(String[] args) {
 
  try {
 
   Workbook wb = new XSSFWorkbook();
   Sheet sheet = wb.createSheet("My Sample Excel");
 
   //FileInputStream obtains input bytes from the image file
   InputStream inputStream = new FileInputStream("../java/img.jpeg");
   //Get the contents of an InputStream as a byte[].
   byte[] bytes = IOUtils.toByteArray(inputStream);
   //Adds a picture to the workbook
   int pictureIdx = wb.addPicture(bytes, Workbook.PICTURE_TYPE_PNG);
   //close the input stream
   inputStream.close();
 
   //Returns an object that handles instantiating concrete classes
   CreationHelper helper = wb.getCreationHelper();
 
   //Creates the top-level drawing patriarch.
   Drawing drawing = sheet.createDrawingPatriarch();
 
   //Create an anchor that is attached to the worksheet
   ClientAnchor anchor = helper.createClientAnchor();
   //set top-left corner for the image
   anchor.setCol1(1);
   anchor.setRow1(2);
 
   //Creates a picture
   Picture pict = drawing.createPicture(anchor, pictureIdx);
   //Reset the image to the original size
   pict.resize();
 
   //Write the Excel file
   FileOutputStream fileOut = null;
   fileOut = new FileOutputStream("../java/myFile.xlsx");
   wb.write(fileOut);
   fileOut.close();
 
  }
  catch (Exception e) {
   System.out.println(e);
  }
 
 
 }
 
}