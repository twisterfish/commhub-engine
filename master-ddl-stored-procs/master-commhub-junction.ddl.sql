
-- Host:    Database: commhub_junction
-- ------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `app_sp_errors`
--

DROP TABLE IF EXISTS `app_sp_errors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `app_sp_errors` (
  `error_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `error_description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`error_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_sp_errors`
--

LOCK TABLES `app_sp_errors` WRITE;
/*!40000 ALTER TABLE `app_sp_errors` DISABLE KEYS */;
/*!40000 ALTER TABLE `app_sp_errors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `device_location`
--

DROP TABLE IF EXISTS `device_location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `device_location` (
  `location_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `device_id` varchar(32) COLLATE utf8_bin NOT NULL DEFAULT '',
  `latitude` float(10,6) NOT NULL,
  `longitude` float(10,6) NOT NULL,
  `time_of_pulse` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`location_id`),
  KEY `grouping_comb1` (`latitude`,`longitude`),
  KEY `this_device_idx` (`device_id`),
  KEY `timestamp_idx` (`time_of_pulse`),
  CONSTRAINT `device_location_ibfk_1` FOREIGN KEY (`device_id`) REFERENCES `end_user_device` (`device_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `device_location`
--

LOCK TABLES `device_location` WRITE;
/*!40000 ALTER TABLE `device_location` DISABLE KEYS */;
/*!40000 ALTER TABLE `device_location` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user`
--

DROP TABLE IF EXISTS `end_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user` (
  `end_user_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `end_user_token` varchar(32) NOT NULL,
  `emid` varchar(255) NOT NULL,
  `pwid` varchar(128) NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_id`),
  UNIQUE KEY `end_user_token` (`end_user_token`),
  UNIQUE KEY `emid` (`emid`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user`
--

LOCK TABLES `end_user` WRITE;
/*!40000 ALTER TABLE `end_user` DISABLE KEYS */;
INSERT INTO `end_user` VALUES (1,'2325418da6fb11e9a58342010a8e0121','johnny.testing@commhubstuff.com','179ad45c6ce2cb97cf1029e212046e81','2012-07-15 12:21:56'),(2,'8a985cfe86ce11e9bc42526af7764f64','jane.testing@commhubstuff.com','179ad45c6ce2cb97cf1029e212046e81','2012-01-03 15:24:41'),(3,'db2e3ddb21b411ea95e10ece0304bc53','jane.testing@gmail.com','179ad45c6ce2cb97cf1029e212046e81','2012-12-18 16:38:44'),(4,'b519c64e21be11eaa2a812500a379b47','','b0dbdc9c67af79396196d55f3eaafd85','2012-12-18 17:49:15'),(5,'96430edd21c511eab91612529ea83c2b','joe.shmoe@commhubstuff.com','179ad45c6ce2cb97cf1029e212046e81','2012-12-18 18:38:30'),(6,'bc9dbc8d227e11ea9d3712cd07fb5023','jane.testing4@gmail.com','179ad45c6ce2cb97cf1029e212046e81','2012-12-19 16:43:51'),(7,'267daa43227f11ea9d3712cd07fb5023','jane.testing5@gmail.com','179ad45c6ce2cb97cf1029e212046e81','2012-12-19 16:46:49'),(8,'ec2ce9e025a111ea90e706c4b7043921','jane.testing99@gmail.com','179ad45c6ce2cb97cf1029e212046e81','2012-12-23 16:33:17');
/*!40000 ALTER TABLE `end_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_device`
--

DROP TABLE IF EXISTS `end_user_device`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_device` (
  `end_user_device_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `end_user_id` int(11) unsigned NOT NULL,
  `device_id` varchar(32) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `date_issued` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_device_id`),
  KEY `end_user_id` (`end_user_id`),
  KEY `device_id` (`device_id`),
  CONSTRAINT `end_user_device_user_id_ibfk_1` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_device`
--

LOCK TABLES `end_user_device` WRITE;
/*!40000 ALTER TABLE `end_user_device` DISABLE KEYS */;
/*!40000 ALTER TABLE `end_user_device` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_group`
--

DROP TABLE IF EXISTS `end_user_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_group` (
  `end_user_group_id` int(11) unsigned NOT NULL,
  `end_user_id` int(11) unsigned NOT NULL,
  `group_name` varchar(64) DEFAULT NULL,
  `group_desc` varchar(255) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_group_id`),
  KEY `end_user_id` (`end_user_id`),
  CONSTRAINT `end_user_group_ibfk_1` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_group`
--

LOCK TABLES `end_user_group` WRITE;
/*!40000 ALTER TABLE `end_user_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `end_user_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_order`
--

DROP TABLE IF EXISTS `end_user_order`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_order` (
  `end_user_order_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `end_user_id` int(11) unsigned NOT NULL,
  `end_user_ship_to_id` int(11) unsigned NOT NULL,
  `work_ticket_id` int(11) unsigned DEFAULT NULL,
  `purchase_order_id` varchar(32) DEFAULT NULL,
  `vendor_invoice_id` varchar(32) DEFAULT NULL,
  `order_date` datetime DEFAULT NULL,
  `paid_date` datetime DEFAULT NULL,
  `vendor_terms` varchar(32) DEFAULT NULL,
  `total_item_count` int(11) unsigned DEFAULT NULL,
  `total_shipping_cost` decimal(8,2) DEFAULT NULL,
  `total_shipping_weight` decimal(8,2) DEFAULT NULL,
  `total_taxes` decimal(8,2) DEFAULT NULL,
  `grand_total` decimal(8,2) DEFAULT NULL,
  `drop_ship_date` datetime DEFAULT NULL,
  `shipped_date` datetime DEFAULT NULL,
  `date_received` datetime DEFAULT NULL,
  `shipped_method` tinyint(4) DEFAULT NULL,
  `shipped_carrier` varchar(32) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_order_id`),
  KEY `end_user_id` (`end_user_id`),
  KEY `order_date` (`order_date`),
  KEY `end_user_order_ibfk_2` (`end_user_ship_to_id`),
  CONSTRAINT `end_user_order_ibfk_1` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `end_user_order_ibfk_2` FOREIGN KEY (`end_user_ship_to_id`) REFERENCES `end_user_ship_to` (`end_user_ship_to_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_order`
--

LOCK TABLES `end_user_order` WRITE;
/*!40000 ALTER TABLE `end_user_order` DISABLE KEYS */;
/*!40000 ALTER TABLE `end_user_order` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_order_item`
--

DROP TABLE IF EXISTS `end_user_order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_order_item` (
  `end_user_order_item_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `end_user_order_id` int(11) unsigned NOT NULL,
  `real_space_inventory_id` int(11) unsigned NOT NULL DEFAULT '0',
  `product_id` int(11) unsigned NOT NULL,
  `product_id_qty` int(11) unsigned NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_order_item_id`),
  KEY `product_id` (`product_id`),
  KEY `end_user_order_id` (`end_user_order_id`),
  CONSTRAINT `end_user_order_item_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`),
  CONSTRAINT `end_user_order_item_ibfk_2` FOREIGN KEY (`end_user_order_id`) REFERENCES `end_user_order` (`end_user_order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_order_item`
--

LOCK TABLES `end_user_order_item` WRITE;
/*!40000 ALTER TABLE `end_user_order_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `end_user_order_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_profile`
--

DROP TABLE IF EXISTS `end_user_profile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_profile` (
  `end_user_profile_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `end_user_id` int(11) unsigned NOT NULL,
  `ein_tax_id` varchar(255) DEFAULT NULL,
  `ssn_tax_id` varchar(255) DEFAULT NULL,
  `last_name` varchar(255) DEFAULT NULL,
  `middle_name` varchar(255) DEFAULT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `address1` varchar(255) DEFAULT NULL,
  `address2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `province_state` varchar(32) DEFAULT NULL,
  `zip_postal_code` varchar(16) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `country_code` varchar(32) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_profile_id`),
  UNIQUE KEY `end_user_id` (`end_user_id`),
  KEY `ein_tax_id` (`ein_tax_id`),
  KEY `ssn_tax_id` (`ssn_tax_id`),
  KEY `last_name` (`last_name`),
  KEY `zip_postal_code` (`zip_postal_code`),
  KEY `full_name` (`first_name`,`last_name`),
  KEY `email` (`email`),
  CONSTRAINT `end_user_profile_id_ibfk_1` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_profile`
--

LOCK TABLES `end_user_profile` WRITE;
/*!40000 ALTER TABLE `end_user_profile` DISABLE KEYS */;
INSERT INTO `end_user_profile` VALUES (1,1,'','','mickelson',NULL,'Johnny',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2012-01-03 15:24:41'),(2,2,'','','Testing',NULL,'Joleen',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2012-01-03 15:24:41'),(3,3,'','','Testing','','Joleen','','','','','','jane.testing@gmail.com','1239958033','','2012-12-18 16:38:44'),(4,4,'','','Munive','','AJ','','','','','','','3057109535','','2012-12-18 17:49:15'),(5,5,'','','Money','','J','','','','','','joe.shmoe@commhubstuff.com','1231231234','','2012-12-18 18:38:30'),(6,6,'','','Testing','','Joleen','','','','','','jane.testing4@gmail.com','1239958033','','2012-12-19 16:43:51'),(7,7,'','','Testing','','Joleen','','','','','','jane.testing5@gmail.com','1239958033','','2012-12-19 16:46:49'),(8,8,'','','Testing2','','Joleen2','','','','','','jane.testing99@gmail.com','9999999999','','2012-12-23 16:33:17');
/*!40000 ALTER TABLE `end_user_profile` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_role`
--

DROP TABLE IF EXISTS `end_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_role` (
  `end_user_role_id` int(11) unsigned NOT NULL,
  `end_user_id` int(11) unsigned NOT NULL,
  `role_name` varchar(64) DEFAULT NULL,
  `role_desc` varchar(255) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_role_id`),
  KEY `end_user_id` (`end_user_id`),
  CONSTRAINT `end_user_role_ibfk_1` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_role`
--

LOCK TABLES `end_user_role` WRITE;
/*!40000 ALTER TABLE `end_user_role` DISABLE KEYS */;
/*!40000 ALTER TABLE `end_user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `end_user_ship_to`
--

DROP TABLE IF EXISTS `end_user_ship_to`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `end_user_ship_to` (
  `end_user_ship_to_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `end_user_id` int(11) unsigned NOT NULL DEFAULT '1',
  `ship_to_first_name` varchar(255) DEFAULT NULL,
  `ship_to_last_name` varchar(255) DEFAULT NULL,
  `ship_to_attention` varchar(255) DEFAULT NULL,
  `ship_to_address1` varchar(255) DEFAULT NULL,
  `ship_to_address2` varchar(255) DEFAULT NULL,
  `ship_to_city` varchar(255) DEFAULT NULL,
  `ship_to_state` varchar(2) DEFAULT NULL,
  `ship_to_zip_postal_code` varchar(16) DEFAULT NULL,
  `ship_to_country_code` varchar(32) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`end_user_ship_to_id`),
  KEY `end_user_ship_to_ibfk_1` (`end_user_id`),
  CONSTRAINT `end_user_ship_to_ibfk_2` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `end_user_ship_to`
--

LOCK TABLES `end_user_ship_to` WRITE;
/*!40000 ALTER TABLE `end_user_ship_to` DISABLE KEYS */;
/*!40000 ALTER TABLE `end_user_ship_to` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invite_status`
--

DROP TABLE IF EXISTS `invite_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invite_status` (
  `invite_status_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`invite_status_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invite_status`
--

LOCK TABLES `invite_status` WRITE;
/*!40000 ALTER TABLE `invite_status` DISABLE KEYS */;
INSERT INTO `invite_status` VALUES (1,'PENDING','2012-07-05 16:29:07'),(2,'ACCEPTED','2012-11-08 16:06:43'),(3,'REJECTED','2012-11-08 16:06:43'),(4,'REVOKED','2012-11-08 16:06:43');
/*!40000 ALTER TABLE `invite_status` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invite_to_workspace`
--

DROP TABLE IF EXISTS `invite_to_workspace`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invite_to_workspace` (
  `invite_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `invite_token` varchar(32) NOT NULL,
  `end_user_id` int(11) unsigned NOT NULL,
  `target_user_id` int(11) unsigned DEFAULT NULL,
  `target_ws_id` int(11) unsigned NOT NULL,
  `target_ws_perm_id` int(11) unsigned NOT NULL,
  `invite_status_id` int(11) unsigned NOT NULL,
  `target_email` varchar(255) NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`invite_id`),
  UNIQUE KEY `invite_token` (`invite_token`),
  KEY `end_user_id` (`end_user_id`),
  KEY `target_user_id` (`target_user_id`),
  KEY `target_ws_id` (`target_ws_id`),
  KEY `target_email` (`target_email`),
  KEY `invite_status_id` (`invite_status_id`),
  KEY `target_ws_perm_id` (`target_ws_perm_id`),
  CONSTRAINT `invite_to_ws_ibfk_1` FOREIGN KEY (`end_user_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `invite_to_ws_ibfk_2` FOREIGN KEY (`target_user_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `invite_to_ws_ibfk_3` FOREIGN KEY (`target_ws_id`) REFERENCES `workspace` (`workspace_id`),
  CONSTRAINT `invite_to_ws_ibfk_4` FOREIGN KEY (`invite_status_id`) REFERENCES `invite_status` (`invite_status_id`),
  CONSTRAINT `invite_to_ws_ibfk_5` FOREIGN KEY (`target_ws_perm_id`) REFERENCES `workspace_permission` (`workspace_permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invite_to_workspace`
--

LOCK TABLES `invite_to_workspace` WRITE;
/*!40000 ALTER TABLE `invite_to_workspace` DISABLE KEYS */;
INSERT INTO `invite_to_workspace` VALUES (1,'5525e852227c11ea9d3712cd07fb5023',1,NULL,1,400,1,'jonparse@email.com','2012-12-19 16:26:38'),(2,'f5f2d86822ab11eaa3440608871fe417',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-19 22:07:34'),(3,'35a9603a22af11eaa3440608871fe417',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-19 22:30:50'),(4,'25938467235c11ea971912facf5bc3d1',1,NULL,1,400,1,'joel.shmoze@outlook.com','2012-12-20 19:08:46'),(5,'531d03c5235c11ea971912facf5bc3d1',1,NULL,1,400,1,'joel.shmoze@commhubstuff.com','2012-12-20 19:10:02'),(6,'8c20842e235c11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 19:11:38'),(7,'62725a36235d11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 19:17:37'),(8,'7debb638236111ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 19:47:01'),(9,'e18bb2bd236211ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 19:56:58'),(10,'ee23af88236311ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:04:29'),(11,'d5c6697c236511ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:18:07'),(12,'21fc4641236611ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:20:15'),(13,'41801d76236611ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:21:08'),(14,'466580d3236611ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:21:16'),(15,'568f377a236611ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:21:43'),(16,'8cf9974d236611ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:23:14'),(17,'17eef77a236711ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:27:07'),(18,'a9a7e6f7236711ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:31:12'),(19,'317cdc09236a11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:49:19'),(20,'3bfa96ef236a11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:49:36'),(21,'bae232fb236a11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:53:09'),(22,'c353d0c4236a11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:53:23'),(23,'fe1d6fa7236a11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:55:02'),(24,'6f0a4b48236b11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 20:58:11'),(25,'1dc712f0236c11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 21:03:05'),(26,'5d866a51236c11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 21:04:52'),(27,'a843f240236c11ea971912facf5bc3d1',1,NULL,1,400,1,'johnny.testing@commhubstuff.com','2012-12-20 21:06:57'),(28,'058ade55236d11ea971912facf5bc3d1',1,3,1,400,2,'jane.testing@commhubstuff.com','2012-12-20 21:10:16'),(29,'12f177b7236d11ea971912facf5bc3d1',1,NULL,1,400,1,'joe.shmoe@commhubstuff.com','2012-12-20 21:09:56'),(30,'e3ca118c25a011ea8f830693478076a1',3,NULL,4,400,1,'jane.testing@gmail.com','2012-12-23 16:25:53'),(31,'4c6232fe25a211ea90e706c4b7043921',3,8,4,400,2,'jane.testing@gmail.com','2012-12-23 16:36:21');
/*!40000 ALTER TABLE `invite_to_workspace` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perishables`
--

DROP TABLE IF EXISTS `perishables`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perishables` (
  `perishables_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `perishables_type_id` int(11) unsigned DEFAULT '1',
  `vendor_id` int(11) unsigned NOT NULL DEFAULT '1',
  `perishables_unit_sold_id` int(11) unsigned DEFAULT '1',
  `upc` int(11) unsigned DEFAULT '1',
  `sku` varchar(32) DEFAULT '1',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`perishables_id`),
  KEY `vendor_id` (`vendor_id`),
  KEY `upc` (`upc`),
  KEY `sku` (`sku`),
  KEY `perishables_type_id` (`perishables_type_id`),
  KEY `perishables_unit_sold_id` (`perishables_unit_sold_id`),
  CONSTRAINT `perishables_ibfk_1` FOREIGN KEY (`vendor_id`) REFERENCES `vendor` (`vendor_id`),
  CONSTRAINT `perishables_ibfk_2` FOREIGN KEY (`perishables_type_id`) REFERENCES `perishables_type` (`perishables_type_id`),
  CONSTRAINT `perishables_ibfk_3` FOREIGN KEY (`perishables_unit_sold_id`) REFERENCES `perishables_unit_sold` (`perishables_unit_sold_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perishables`
--

LOCK TABLES `perishables` WRITE;
/*!40000 ALTER TABLE `perishables` DISABLE KEYS */;
/*!40000 ALTER TABLE `perishables` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perishables_component_list`
--

DROP TABLE IF EXISTS `perishables_component_list`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perishables_component_list` (
  `perishables_component_list_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned DEFAULT '1',
  `perishables_id` int(11) unsigned DEFAULT '1',
  `perishables_uom_id` int(11) unsigned DEFAULT '1',
  `perishable_qty` int(11) unsigned DEFAULT '1',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`perishables_component_list_id`),
  KEY `product_id` (`product_id`),
  KEY `perishables_id` (`perishables_id`),
  KEY `perishables_uom_id` (`perishables_uom_id`),
  CONSTRAINT `perishables_component_list_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`),
  CONSTRAINT `perishables_component_list_ibfk_2` FOREIGN KEY (`perishables_id`) REFERENCES `perishables` (`perishables_id`),
  CONSTRAINT `perishables_component_list_ibfk_3` FOREIGN KEY (`perishables_uom_id`) REFERENCES `perishables_unit_of_measure` (`perishables_unit_of_measure_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perishables_component_list`
--

LOCK TABLES `perishables_component_list` WRITE;
/*!40000 ALTER TABLE `perishables_component_list` DISABLE KEYS */;
/*!40000 ALTER TABLE `perishables_component_list` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perishables_cost_history`
--

DROP TABLE IF EXISTS `perishables_cost_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perishables_cost_history` (
  `perishables_cost_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `perishables_id` int(11) unsigned NOT NULL,
  `cost` decimal(8,2) DEFAULT NULL,
  `effective_date` datetime DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`perishables_cost_id`),
  KEY `perishables_id` (`perishables_id`),
  KEY `effective_date` (`effective_date`),
  CONSTRAINT `perishables_cost_history_ibfk_1` FOREIGN KEY (`perishables_id`) REFERENCES `perishables` (`perishables_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perishables_cost_history`
--

LOCK TABLES `perishables_cost_history` WRITE;
/*!40000 ALTER TABLE `perishables_cost_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `perishables_cost_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perishables_type`
--

DROP TABLE IF EXISTS `perishables_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perishables_type` (
  `perishables_type_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`perishables_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perishables_type`
--

LOCK TABLES `perishables_type` WRITE;
/*!40000 ALTER TABLE `perishables_type` DISABLE KEYS */;
INSERT INTO `perishables_type` VALUES (1,'fish','2012-10-02 15:25:37'),(2,'beef','2012-10-02 15:25:43'),(3,'chicken','2012-10-02 15:25:55'),(4,'pork','2012-10-02 15:26:27'),(5,'fowl','2012-10-02 15:27:07'),(6,'exotic meats','2012-10-02 15:27:15'),(7,'fruit','2012-10-02 15:27:26'),(8,'vegetables','2012-10-02 15:27:45'),(9,'grains','2012-10-02 15:28:01'),(10,'dairy','2012-10-02 15:30:07'),(11,'carbonated beverage','2012-10-02 15:30:35'),(12,'carbonated alcoholic beverage','2012-10-02 15:31:14'),(13,'wine','2012-10-02 15:32:25'),(14,'liquor','2012-10-02 15:38:40');
/*!40000 ALTER TABLE `perishables_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perishables_unit_of_measure`
--

DROP TABLE IF EXISTS `perishables_unit_of_measure`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perishables_unit_of_measure` (
  `perishables_unit_of_measure_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`perishables_unit_of_measure_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perishables_unit_of_measure`
--

LOCK TABLES `perishables_unit_of_measure` WRITE;
/*!40000 ALTER TABLE `perishables_unit_of_measure` DISABLE KEYS */;
INSERT INTO `perishables_unit_of_measure` VALUES (1,'teaspoon','2012-10-03 14:25:57'),(2,'dessertspoon','2012-10-03 14:26:10'),(3,'tablespoon','2012-10-03 14:26:21'),(4,'fluid ounce','2012-10-03 14:26:33'),(5,'cup','2012-10-03 14:27:01'),(6,'pint','2012-10-03 14:27:07'),(7,'quart','2012-10-03 14:27:15'),(8,'gallon','2012-10-03 14:27:37'),(9,'liquid barrel','2012-10-03 14:28:39'),(10,'dry ounce','2012-10-03 15:21:11'),(11,'pound','2012-10-03 15:21:37'),(12,'ton','2012-10-03 15:21:44'),(13,'milliliter','2012-10-03 15:23:02'),(14,'liter','2012-10-03 15:23:10'),(15,'kiloliter','2012-10-03 15:23:41'),(16,'milligram','2012-10-03 15:23:49'),(17,'gram','2012-10-03 15:23:55'),(18,'kilogram','2012-10-03 15:24:24'),(19,'metric tonne','2012-10-03 15:24:56'),(20,'peck','2012-10-03 15:26:15'),(21,'bushel','2012-10-03 15:26:22');
/*!40000 ALTER TABLE `perishables_unit_of_measure` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perishables_unit_sold`
--

DROP TABLE IF EXISTS `perishables_unit_sold`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perishables_unit_sold` (
  `perishables_unit_sold_id` int(11) unsigned NOT NULL,
  `qty_in_unit` int(11) unsigned NOT NULL DEFAULT '1',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`perishables_unit_sold_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perishables_unit_sold`
--

LOCK TABLES `perishables_unit_sold` WRITE;
/*!40000 ALTER TABLE `perishables_unit_sold` DISABLE KEYS */;
/*!40000 ALTER TABLE `perishables_unit_sold` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product`
--

DROP TABLE IF EXISTS `product`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product` (
  `product_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` int(11) unsigned DEFAULT '1',
  `vendor_id` int(11) unsigned NOT NULL DEFAULT '1',
  `product_unit_sold_id` int(11) unsigned DEFAULT '1',
  `upc` varchar(64) DEFAULT '1',
  `sku` varchar(64) DEFAULT '1',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_id`),
  KEY `vendor_id` (`vendor_id`),
  KEY `upc` (`upc`),
  KEY `sku` (`sku`),
  KEY `workspace_id` (`workspace_id`),
  KEY `product_unit_sold_id` (`product_unit_sold_id`),
  CONSTRAINT `product_ibfk_1` FOREIGN KEY (`vendor_id`) REFERENCES `vendor` (`vendor_id`),
  CONSTRAINT `product_ibfk_2` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`workspace_id`),
  CONSTRAINT `product_ibfk_3` FOREIGN KEY (`product_unit_sold_id`) REFERENCES `product_unit_sold` (`product_unit_sold_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product`
--

LOCK TABLES `product` WRITE;
/*!40000 ALTER TABLE `product` DISABLE KEYS */;
INSERT INTO `product` VALUES (1,1,1,1,'983459345934345','234567Y','cool ass widget','2012-01-31 15:38:27'),(2,1,1,1,'1','1','1','2012-04-09 19:02:28'),(3,4,1,1,'UPC1','SKU1','1','2012-12-23 20:26:24');
/*!40000 ALTER TABLE `product` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_asset`
--

DROP TABLE IF EXISTS `product_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_asset` (
  `product_asset_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned NOT NULL DEFAULT '1',
  `uploader_id` int(11) unsigned NOT NULL DEFAULT '1',
  `asset_name` varchar(32) NOT NULL DEFAULT '0',
  `content_type` varchar(32) NOT NULL DEFAULT '0',
  `asset_size_bytes` int(11) DEFAULT '0',
  `asset_size_height` int(11) DEFAULT '0',
  `asset_size_width` int(11) DEFAULT '0',
  `asset_upload_path` text,
  `asset_description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_asset_id`),
  KEY `asset_name` (`asset_name`),
  KEY `product_id` (`product_id`),
  CONSTRAINT `product_asset_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1045 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_asset`
--

LOCK TABLES `product_asset` WRITE;
/*!40000 ALTER TABLE `product_asset` DISABLE KEYS */;
INSERT INTO `product_asset` VALUES (1,1,1,'1','1',1,1,1,'1','1','2012-02-27 16:48:40'),(2,1,1,'1','1',1,1,1,'1','1','2012-02-27 16:49:27'),(3,1,1,'1','1',1,1,1,'1','1','2012-02-27 16:50:27'),(4,1,1,'1','1',1,1,1,'1','1','2012-02-27 16:53:00'),(5,1,1,'1','1',1,1,1,'1','1','2012-02-27 16:54:29'),(6,1,1,'1','1',1,1,1,'1','1','2012-02-27 16:56:29'),(7,1,1,'1','1',1,1,1,'1','1','2012-02-27 19:39:00'),(8,1,1,'1','1',1,1,1,'1','1','2012-02-27 19:39:53'),(9,1,1,'1','1',1,1,1,'1','1','2012-02-27 19:40:02'),(10,1,1,'1','1',1,1,1,'1','1','2012-02-27 19:54:57'),(11,1,1,'1','1',1,1,1,'1','1','2012-03-01 16:39:38'),(12,1,1,'1','1',1,1,1,'1','1','2012-03-05 17:08:27'),(13,1,1,'1','1',1,1,1,'1','1','2012-03-06 20:29:12'),(14,1,1,'1','1',1,1,1,'1','1','2012-03-06 20:51:40'),(15,1,1,'1','1',1,1,1,'1','1','2012-03-11 15:34:41'),(16,1,1,'1','1',1,1,1,'1','1','2012-03-15 13:30:43'),(17,1,1,'1','1',1,1,1,'1','1','2012-03-15 16:03:17'),(18,1,1,'1','1',1,1,1,'1','1','2012-03-25 17:53:46'),(19,1,1,'1','1',1,1,1,'1','1','2012-03-25 18:06:43'),(20,1,1,'1','1',1,1,1,'1','1','2012-03-25 18:07:30'),(21,1,1,'1','1',1,1,1,'1','1','2012-03-25 18:09:56'),(22,1,1,'1','1',1,1,1,'1','1','2012-03-25 20:34:11'),(23,1,1,'1','1',1,1,1,'1','1','2012-03-25 20:34:57'),(24,1,1,'1','1',1,1,1,'1','1','2012-03-25 20:35:41'),(25,1,1,'1','1',1,1,1,'1','1','2012-04-02 16:31:30'),(26,1,1,'1','1',1,1,1,'1','1','2012-04-04 19:20:57'),(27,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:14:24'),(28,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:15:37'),(29,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:16:39'),(30,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:20:05'),(31,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:20:46'),(32,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:37:25'),(33,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:41:20'),(34,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:44:59'),(35,1,1,'1','1',1,1,1,'1','1','2012-04-09 15:49:44'),(36,1,1,'1','1',1,1,1,'1','1','2012-04-09 17:02:31'),(37,1,1,'1','1',1,1,1,'1','1','2012-04-09 19:02:28'),(38,1,1,'1','1',1,1,1,'1','1','2012-04-09 20:58:45'),(1039,1,1,'1','1',1,1,1,'1','1','2012-07-18 15:35:58'),(1040,1,1,'1','1',1,1,1,'1','1','2012-07-18 16:08:15'),(1041,1,1,'1','1',1,1,1,'1','1','2012-07-18 16:13:00'),(1042,1,1,'1','1',1,1,1,'1','1','2012-07-18 16:15:19'),(1043,1,1,'1','1',1,1,1,'1','1','2012-07-26 13:56:03'),(1044,1,1,'1','1',1,1,1,'1','1','2012-07-26 13:58:12');
/*!40000 ALTER TABLE `product_asset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_cost_history`
--

DROP TABLE IF EXISTS `product_cost_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_cost_history` (
  `product_cost_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned NOT NULL,
  `cost` decimal(8,2) DEFAULT NULL,
  `effective_date` datetime DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_cost_id`),
  KEY `product_id` (`product_id`),
  KEY `effective_date` (`effective_date`),
  CONSTRAINT `product_cost_history_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_cost_history`
--

LOCK TABLES `product_cost_history` WRITE;
/*!40000 ALTER TABLE `product_cost_history` DISABLE KEYS */;
INSERT INTO `product_cost_history` VALUES (1,1,5.00,'1970-01-01 00:00:01','2012-01-31 15:53:30'),(2,1,4.50,'2016-01-01 00:00:01','2012-02-01 20:28:56'),(3,1,1.00,'1970-01-01 02:46:40','2012-04-09 15:44:59'),(4,1,1.00,'1970-01-01 02:46:40','2012-04-09 15:49:44'),(5,1,1.00,'1970-01-01 02:46:40','2012-04-09 17:02:31'),(6,1,1.00,'1970-01-01 02:46:40','2012-04-09 19:02:28'),(7,1,1.00,'1970-01-01 02:46:40','2012-04-09 20:58:45'),(8,1,103.00,'1970-01-13 20:57:02','2012-12-18 20:20:47');
/*!40000 ALTER TABLE `product_cost_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_inventory`
--

DROP TABLE IF EXISTS `product_inventory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_inventory` (
  `product_inventory_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned NOT NULL,
  `total_quantity_on_hand` decimal(8,2) DEFAULT NULL,
  `order_threshold` decimal(8,2) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_inventory_id`),
  UNIQUE KEY `product_id` (`product_id`),
  CONSTRAINT `product_inventory_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_inventory`
--

LOCK TABLES `product_inventory` WRITE;
/*!40000 ALTER TABLE `product_inventory` DISABLE KEYS */;
INSERT INTO `product_inventory` VALUES (1,1,944.00,NULL,'2012-12-19 16:18:51'),(2,3,790.00,NULL,'2012-12-23 20:46:00');
/*!40000 ALTER TABLE `product_inventory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_location`
--

DROP TABLE IF EXISTS `product_location`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_location` (
  `location_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned NOT NULL,
  `latitude` float(10,6) NOT NULL,
  `longitude` float(10,6) NOT NULL,
  `quantity_on_location` decimal(8,2) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`location_id`),
  KEY `grouping_comb1` (`latitude`,`longitude`),
  KEY `product_loc_idx` (`product_id`),
  KEY `timestamp_idx` (`last_updated`),
  CONSTRAINT `product_location_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_location`
--

LOCK TABLES `product_location` WRITE;
/*!40000 ALTER TABLE `product_location` DISABLE KEYS */;
/*!40000 ALTER TABLE `product_location` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_price_history`
--

DROP TABLE IF EXISTS `product_price_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_price_history` (
  `product_price_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int(11) unsigned NOT NULL,
  `price` decimal(8,2) DEFAULT NULL,
  `effective_date` datetime DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_price_id`),
  KEY `product_id` (`product_id`),
  KEY `effective_date` (`effective_date`),
  CONSTRAINT `product_price_history_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_price_history`
--

LOCK TABLES `product_price_history` WRITE;
/*!40000 ALTER TABLE `product_price_history` DISABLE KEYS */;
INSERT INTO `product_price_history` VALUES (1,1,10.00,'1970-01-01 00:00:01','2012-01-31 15:46:04'),(2,2,1.00,'1970-01-01 02:46:40','2012-04-09 15:41:20'),(3,1,105.00,'1972-11-20 12:57:02','2012-12-18 20:26:39');
/*!40000 ALTER TABLE `product_price_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_type`
--

DROP TABLE IF EXISTS `product_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_type` (
  `product_type_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_type`
--

LOCK TABLES `product_type` WRITE;
/*!40000 ALTER TABLE `product_type` DISABLE KEYS */;
INSERT INTO `product_type` VALUES (1,'commercial hardware','2012-09-26 18:58:43'),(2,'commercial plumbing','2012-09-26 18:58:55'),(3,'commercial electrical','2012-09-26 18:59:05'),(4,'wholesale perishables','2012-09-26 18:59:23'),(5,'wholesale beverage','2012-09-26 18:59:45'),(6,'individual beverage','2012-09-26 19:00:37'),(7,'individual meal','2012-09-26 19:01:22'),(8,'lumber','2012-09-26 19:01:52'),(9,'Toys','2012-01-31 15:40:41');
/*!40000 ALTER TABLE `product_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `product_unit_sold`
--

DROP TABLE IF EXISTS `product_unit_sold`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `product_unit_sold` (
  `product_unit_sold_id` int(11) unsigned NOT NULL,
  `qty_in_unit` float(10,4) NOT NULL DEFAULT '1.0000',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`product_unit_sold_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `product_unit_sold`
--

LOCK TABLES `product_unit_sold` WRITE;
/*!40000 ALTER TABLE `product_unit_sold` DISABLE KEYS */;
INSERT INTO `product_unit_sold` VALUES (1,1.0000,'each','2012-09-21 14:59:40'),(2,5.0000,'plastic pack','2012-09-21 15:00:59'),(3,25.0000,'box','2012-09-21 15:01:31'),(4,100.0000,'carton','2012-09-21 15:01:49'),(5,10000.0000,'crate','2012-09-21 15:02:49'),(6,1000000.0000,'freight container','2012-09-21 15:04:05');
/*!40000 ALTER TABLE `product_unit_sold` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_property`
--

DROP TABLE IF EXISTS `real_property`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_property` (
  `real_property_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `real_property_type_id` int(11) unsigned NOT NULL DEFAULT '1',
  `real_property_name` varchar(255) DEFAULT NULL,
  `address1` varchar(255) DEFAULT NULL,
  `address2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `province_state` varchar(32) DEFAULT NULL,
  `zip_postal_code` varchar(16) DEFAULT NULL,
  `country_code` varchar(32) DEFAULT NULL,
  `latitude` float(10,6) NOT NULL,
  `longitude` float(10,6) NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_property_id`),
  KEY `real_property_type_id` (`real_property_type_id`),
  CONSTRAINT `real_property_type_ibfk_1` FOREIGN KEY (`real_property_type_id`) REFERENCES `real_property_type` (`real_property_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_property`
--

LOCK TABLES `real_property` WRITE;
/*!40000 ALTER TABLE `real_property` DISABLE KEYS */;
INSERT INTO `real_property` VALUES (1,1,'Ocean Hotel','5260 Breeze Blvd.','','Fort Meyer','VA','33316','US',26.108900,-80.106735,'2012-04-24 20:26:43'),(2,1,'Resort & Spa','1005 Plaza Blvd','','Lake Vista','VA','76901','US',28.377148,-81.506920,'2012-04-24 20:26:43'),(3,1,'Riverwalk','2001 N Ash Dr.','','Tampa','VA','33602','US',27.946451,-82.459564,'2012-04-24 20:26:43'),(4,1,'Canal','1000 Canal St','','New Orleans','LA','70112','US',29.956497,-90.074333,'2012-04-24 20:26:43'),(5,1,'Hotel and Water Resort','70 Internet Dr','','Orlando','VA','32819','US',28.454250,-81.471527,'2012-04-24 20:26:43'),(6,1,'Resort & Marina','450 OverHwy','','Marathon','VA','33050','US',24.714798,-81.083946,'2012-04-24 20:26:43'),(7,1,'commhub Office','910 SE 17th St','','Fort Meyer','VA','33316','US',26.100100,-80.133202,'2012-05-16 18:04:41');
/*!40000 ALTER TABLE `real_property` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_property_asset`
--

DROP TABLE IF EXISTS `real_property_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_property_asset` (
  `real_property_asset_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `real_property_id` int(11) unsigned NOT NULL DEFAULT '1',
  `uploader_id` int(11) unsigned NOT NULL DEFAULT '1',
  `asset_name` varchar(32) NOT NULL DEFAULT '0',
  `content_type` varchar(32) NOT NULL DEFAULT '0',
  `asset_size_bytes` int(11) DEFAULT '0',
  `asset_size_height` int(11) DEFAULT '0',
  `asset_size_width` int(11) DEFAULT '0',
  `asset_upload_path` text,
  `asset_description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_property_asset_id`),
  KEY `asset_name` (`asset_name`),
  KEY `real_property_id` (`real_property_id`),
  CONSTRAINT `real_property_asset_ibfk_1` FOREIGN KEY (`real_property_id`) REFERENCES `real_property` (`real_property_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_property_asset`
--

LOCK TABLES `real_property_asset` WRITE;
/*!40000 ALTER TABLE `real_property_asset` DISABLE KEYS */;
/*!40000 ALTER TABLE `real_property_asset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_property_fixture`
--

DROP TABLE IF EXISTS `real_property_fixture`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_property_fixture` (
  `real_property_fixture_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `rpf_token` varchar(32) NOT NULL,
  `real_property_id` int(11) unsigned NOT NULL DEFAULT '1',
  `upc` varchar(64) DEFAULT '1',
  `sku` varchar(64) DEFAULT '1',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_property_fixture_id`),
  UNIQUE KEY `rpf_token` (`rpf_token`),
  KEY `upc` (`upc`),
  KEY `sku` (`sku`),
  KEY `real_property_fixture_ibfk_1` (`real_property_id`),
  CONSTRAINT `real_property_fixture_ibfk_1` FOREIGN KEY (`real_property_id`) REFERENCES `real_property` (`real_property_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_property_fixture`
--

LOCK TABLES `real_property_fixture` WRITE;
/*!40000 ALTER TABLE `real_property_fixture` DISABLE KEYS */;
/*!40000 ALTER TABLE `real_property_fixture` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_property_type`
--

DROP TABLE IF EXISTS `real_property_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_property_type` (
  `real_property_type_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_property_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_property_type`
--

LOCK TABLES `real_property_type` WRITE;
/*!40000 ALTER TABLE `real_property_type` DISABLE KEYS */;
INSERT INTO `real_property_type` VALUES (1,'comercial property','2012-10-02 15:17:58'),(2,'multi family residence','2012-10-02 15:18:18'),(3,'single family residence','2012-10-02 15:18:29'),(4,'commercial motel','2012-10-02 15:18:58'),(5,'commercial resort','2012-10-02 15:19:07');
/*!40000 ALTER TABLE `real_property_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_space_inventory`
--

DROP TABLE IF EXISTS `real_space_inventory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_space_inventory` (
  `real_space_inventory_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `real_property_id` int(11) unsigned NOT NULL DEFAULT '1',
  `real_space_inventory_type_id` int(11) unsigned NOT NULL DEFAULT '1',
  `real_space_inventory_status_id` int(11) unsigned NOT NULL DEFAULT '1',
  `space_length` decimal(8,2) DEFAULT NULL,
  `space_width` decimal(8,2) DEFAULT NULL,
  `space_height` decimal(8,2) DEFAULT NULL,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_space_inventory_id`),
  KEY `real_space_inventory_ibfk_10` (`real_property_id`),
  KEY `real_space_inventory_ibfk_2` (`real_space_inventory_type_id`),
  KEY `real_space_inventory_ibfk_3` (`real_space_inventory_status_id`),
  CONSTRAINT `real_space_inventory_ibfk_1` FOREIGN KEY (`real_property_id`) REFERENCES `real_property` (`real_property_id`),
  CONSTRAINT `real_space_inventory_ibfk_2` FOREIGN KEY (`real_space_inventory_type_id`) REFERENCES `real_space_inventory_type` (`real_space_inventory_type_id`),
  CONSTRAINT `real_space_inventory_ibfk_3` FOREIGN KEY (`real_space_inventory_status_id`) REFERENCES `real_space_inventory_status` (`real_space_inventory_status_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_space_inventory`
--

LOCK TABLES `real_space_inventory` WRITE;
/*!40000 ALTER TABLE `real_space_inventory` DISABLE KEYS */;
/*!40000 ALTER TABLE `real_space_inventory` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_space_inventory_feature`
--

DROP TABLE IF EXISTS `real_space_inventory_feature`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_space_inventory_feature` (
  `real_space_inventory_feature_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `real_space_inventory_id` int(11) unsigned NOT NULL DEFAULT '1',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_space_inventory_feature_id`),
  KEY `real_space_inventory_id` (`real_space_inventory_id`),
  CONSTRAINT `real_space_inventory_feature_ibfk_1` FOREIGN KEY (`real_space_inventory_id`) REFERENCES `real_space_inventory` (`real_space_inventory_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_space_inventory_feature`
--

LOCK TABLES `real_space_inventory_feature` WRITE;
/*!40000 ALTER TABLE `real_space_inventory_feature` DISABLE KEYS */;
/*!40000 ALTER TABLE `real_space_inventory_feature` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_space_inventory_feedback`
--

DROP TABLE IF EXISTS `real_space_inventory_feedback`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_space_inventory_feedback` (
  `real_space_inventory_feedback_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `real_space_inventory_id` int(11) unsigned NOT NULL DEFAULT '1',
  `feedback` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_space_inventory_feedback_id`),
  KEY `real_space_inventory_feedback_ibfk_1` (`real_space_inventory_id`),
  CONSTRAINT `real_space_inventory_feedback_ibfk_1` FOREIGN KEY (`real_space_inventory_id`) REFERENCES `real_space_inventory` (`real_space_inventory_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_space_inventory_feedback`
--

LOCK TABLES `real_space_inventory_feedback` WRITE;
/*!40000 ALTER TABLE `real_space_inventory_feedback` DISABLE KEYS */;
/*!40000 ALTER TABLE `real_space_inventory_feedback` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_space_inventory_price_history`
--

DROP TABLE IF EXISTS `real_space_inventory_price_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_space_inventory_price_history` (
  `rsi_price_history_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `real_space_inventory_id` int(11) unsigned NOT NULL DEFAULT '1',
  `price` decimal(8,2) DEFAULT NULL,
  `effective_date` datetime DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`rsi_price_history_id`),
  KEY `real_space_inventory_id` (`real_space_inventory_id`),
  KEY `effective_date` (`effective_date`),
  CONSTRAINT `rsi_price_history_ibfk_1` FOREIGN KEY (`real_space_inventory_id`) REFERENCES `real_space_inventory` (`real_space_inventory_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_space_inventory_price_history`
--

LOCK TABLES `real_space_inventory_price_history` WRITE;
/*!40000 ALTER TABLE `real_space_inventory_price_history` DISABLE KEYS */;
/*!40000 ALTER TABLE `real_space_inventory_price_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_space_inventory_status`
--

DROP TABLE IF EXISTS `real_space_inventory_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_space_inventory_status` (
  `real_space_inventory_status_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_space_inventory_status_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_space_inventory_status`
--

LOCK TABLES `real_space_inventory_status` WRITE;
/*!40000 ALTER TABLE `real_space_inventory_status` DISABLE KEYS */;
INSERT INTO `real_space_inventory_status` VALUES (1,'occupied needs repairs','2012-11-08 18:25:51'),(2,'vacant needs repairs','2012-11-08 18:25:51'),(3,'occupied needs cleaning','2012-11-08 18:25:51'),(4,'vacant needs cleaning','2012-11-08 18:25:51'),(5,'ready vacant','2012-11-08 18:25:51'),(6,'occupied','2012-11-08 18:25:51');
/*!40000 ALTER TABLE `real_space_inventory_status` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `real_space_inventory_type`
--

DROP TABLE IF EXISTS `real_space_inventory_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `real_space_inventory_type` (
  `real_space_inventory_type_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`real_space_inventory_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `real_space_inventory_type`
--

LOCK TABLES `real_space_inventory_type` WRITE;
/*!40000 ALTER TABLE `real_space_inventory_type` DISABLE KEYS */;
INSERT INTO `real_space_inventory_type` VALUES (1,'motel room single','2012-10-02 15:20:46'),(2,'motel room suite','2012-10-02 15:20:58'),(3,'resort room suite','2012-10-02 15:21:56'),(4,'resort room single','2012-10-02 15:22:06'),(5,'resort room banquet','2012-10-02 15:22:22'),(6,'resort room conference','2012-10-02 15:22:31');
/*!40000 ALTER TABLE `real_space_inventory_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `session_token`
--

DROP TABLE IF EXISTS `session_token`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `session_token` (
  `token_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `token` varchar(32) NOT NULL,
  `debug_mode` int(11) unsigned NOT NULL DEFAULT '0',
  `sample_mode` int(11) unsigned NOT NULL DEFAULT '0',
  `sample_time_start` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `sample_rate` int(11) unsigned NOT NULL DEFAULT '0',
  `sample_duration` int(11) unsigned NOT NULL DEFAULT '0',
  `label` varchar(128) NOT NULL,
  `lastUpdate` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`token_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `session_token`
--

LOCK TABLES `session_token` WRITE;
/*!40000 ALTER TABLE `session_token` DISABLE KEYS */;
INSERT INTO `session_token` VALUES (1,'37b176076fc74698be5aed02f74cbf15',0,0,1,0,0,'API TOKEN','2012-01-20 21:43:20'),(2,'51a176076fc746988lkP2d092f74c278',0,0,1,0,0,'REGISTRATION TOKEN','2012-05-08 15:53:55'),(3,'62b176076fc746988lkP2d092f74c389',0,0,1,0,0,'DEV MODE TOKEN','2012-09-30 18:14:00');
/*!40000 ALTER TABLE `session_token` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ticket_group`
--

DROP TABLE IF EXISTS `ticket_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ticket_group` (
  `ticket_group_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `ticket_group_creator_id` int(11) unsigned NOT NULL DEFAULT '1',
  `workspace_id` int(11) unsigned NOT NULL DEFAULT '1',
  `ticket_group_name` varchar(64) NOT NULL DEFAULT '0',
  `ticket_group_description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ticket_group_id`),
  KEY `ticket_group_creator_id` (`ticket_group_creator_id`),
  KEY `workspace_id` (`workspace_id`),
  KEY `ticket_group_name` (`ticket_group_name`),
  CONSTRAINT `ticket_group_ibfk_1` FOREIGN KEY (`ticket_group_creator_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `ticket_group_ibfk_3` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`workspace_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ticket_group`
--

LOCK TABLES `ticket_group` WRITE;
/*!40000 ALTER TABLE `ticket_group` DISABLE KEYS */;
INSERT INTO `ticket_group` VALUES (1,1,1,'Default','something','2012-05-20 19:10:59'),(2,2,2,'Default','something','2012-05-07 18:46:57'),(3,3,4,'STAGING','Tickets that are not yet assigned to a group or were rejected by the ticket assignee','2012-12-18 16:38:44'),(4,4,5,'STAGING','Tickets that are not yet assigned to a group or were rejected by the ticket assignee','2012-12-18 17:49:15'),(5,5,6,'STAGING','Tickets that are not yet assigned to a group or were rejected by the ticket assignee','2012-12-18 18:38:30'),(6,6,7,'STAGING','Tickets that are not yet assigned to a group or were rejected by the ticket assignee','2012-12-19 16:43:51'),(7,7,8,'STAGING','Tickets that are not yet assigned to a group or were rejected by the ticket assignee','2012-12-19 16:46:49'),(8,8,9,'STAGING','Tickets that are not yet assigned to a group or were rejected by the ticket assignee','2012-12-23 16:33:17');
/*!40000 ALTER TABLE `ticket_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ticket_group_label`
--

DROP TABLE IF EXISTS `ticket_group_label`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ticket_group_label` (
  `ticket_group_label_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `ticket_group_id` int(11) unsigned NOT NULL DEFAULT '1',
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`ticket_group_label_id`),
  KEY `ticket_group_label_ibfk_1` (`ticket_group_id`),
  CONSTRAINT `ticket_group_label_ibfk_1` FOREIGN KEY (`ticket_group_id`) REFERENCES `ticket_group` (`ticket_group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ticket_group_label`
--

LOCK TABLES `ticket_group_label` WRITE;
/*!40000 ALTER TABLE `ticket_group_label` DISABLE KEYS */;
/*!40000 ALTER TABLE `ticket_group_label` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ticket_rsi_xref`
--

DROP TABLE IF EXISTS `ticket_rsi_xref`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ticket_rsi_xref` (
  `tktrsixrf_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `work_ticket_id` int(11) unsigned NOT NULL,
  `rsi_id` int(11) unsigned NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`tktrsixrf_id`),
  KEY `grouping_comb1` (`work_ticket_id`,`rsi_id`),
  KEY `work_ticket_id` (`work_ticket_id`),
  KEY `rsi_id` (`rsi_id`),
  KEY `last_updated` (`last_updated`),
  CONSTRAINT `wo_rsi_xref_ibfk_1` FOREIGN KEY (`work_ticket_id`) REFERENCES `work_ticket` (`work_ticket_id`),
  CONSTRAINT `wo_rsi_xref_ibfk_2` FOREIGN KEY (`rsi_id`) REFERENCES `real_space_inventory` (`real_space_inventory_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ticket_rsi_xref`
--

LOCK TABLES `ticket_rsi_xref` WRITE;
/*!40000 ALTER TABLE `ticket_rsi_xref` DISABLE KEYS */;
/*!40000 ALTER TABLE `ticket_rsi_xref` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vendor`
--

DROP TABLE IF EXISTS `vendor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vendor` (
  `vendor_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `vendor_type_id` int(11) unsigned NOT NULL DEFAULT '1',
  `vendor_name` varchar(255) DEFAULT NULL,
  `address1` varchar(255) DEFAULT NULL,
  `address2` varchar(255) DEFAULT NULL,
  `city` varchar(255) DEFAULT NULL,
  `province_state` varchar(32) DEFAULT NULL,
  `zip_postal_code` varchar(16) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `country_code` varchar(32) DEFAULT NULL,
  `invoice_to_first_name` varchar(255) DEFAULT NULL,
  `invoice_to_last_name` varchar(255) DEFAULT NULL,
  `invoice_to_attention` varchar(255) DEFAULT NULL,
  `invoice_to_address1` varchar(255) DEFAULT NULL,
  `invoice_to_address2` varchar(255) DEFAULT NULL,
  `invoice_to_city` varchar(255) DEFAULT NULL,
  `invoice_to_state` varchar(2) DEFAULT NULL,
  `invoice_to_zip_postal_code` varchar(16) DEFAULT NULL,
  `invoice_to_country_code` varchar(32) DEFAULT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`vendor_id`),
  KEY `vendor_type_id` (`vendor_type_id`),
  CONSTRAINT `vendor_type_ibfk_1` FOREIGN KEY (`vendor_type_id`) REFERENCES `vendor_type` (`vendor_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vendor`
--

LOCK TABLES `vendor` WRITE;
/*!40000 ALTER TABLE `vendor` DISABLE KEYS */;
INSERT INTO `vendor` VALUES (1,2,'Fast Eddie Goods','910 Se 17th Street','#4','Fort Meyer','VA','33316','johnny.testing@commhub.tech','123-224-5945','1','Johnny','mickelson','JOE SCHMO','910 Se 17th Street','#4','Fort Meyer','VA','33316','1','2012-01-31 19:40:36');
/*!40000 ALTER TABLE `vendor` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `vendor_type`
--

DROP TABLE IF EXISTS `vendor_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `vendor_type` (
  `vendor_type_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`vendor_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `vendor_type`
--

LOCK TABLES `vendor_type` WRITE;
/*!40000 ALTER TABLE `vendor_type` DISABLE KEYS */;
INSERT INTO `vendor_type` VALUES (1,'Construction Contractor','2012-09-20 21:17:46'),(2,'Wholesale Dry Goods','2012-09-21 13:05:23'),(3,'Wholesale Perishables','2012-09-21 13:05:49'),(4,'Retail Dry Goods','2012-09-21 13:06:23'),(5,'Retail Perishables','2012-09-21 13:06:38'),(6,'Commercial Real Estate','2012-09-21 13:11:17');
/*!40000 ALTER TABLE `vendor_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `work_ticket`
--

DROP TABLE IF EXISTS `work_ticket`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `work_ticket` (
  `work_ticket_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `local_guid` varchar(32) NOT NULL DEFAULT '0',
  `ticket_group_id` int(11) unsigned NOT NULL DEFAULT '1',
  `work_ticket_status_id` int(11) unsigned NOT NULL DEFAULT '1',
  `work_ticket_type_id` int(11) unsigned NOT NULL DEFAULT '1',
  `real_property_id` int(11) unsigned NOT NULL DEFAULT '1',
  `creator_user_id` int(11) unsigned NOT NULL DEFAULT '1',
  `assigned_to_user_id` int(11) unsigned NOT NULL DEFAULT '1',
  `assigned_by_user_id` int(11) unsigned NOT NULL DEFAULT '1',
  `running_time` int(11) unsigned NOT NULL DEFAULT '1',
  `time_created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `time_started` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `ts_latitude` float(10,6) DEFAULT '0.000000',
  `ts_longitude` float(10,6) DEFAULT '0.000000',
  `time_finished` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00',
  `tf_latitude` float(10,6) DEFAULT '0.000000',
  `tf_longitude` float(10,6) DEFAULT '0.000000',
  `title` varchar(255) DEFAULT '',
  `description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`work_ticket_id`),
  UNIQUE KEY `local_guid` (`local_guid`),
  KEY `ticket_group_id` (`ticket_group_id`),
  KEY `work_ticket_type_id` (`work_ticket_type_id`),
  KEY `real_property_id` (`real_property_id`),
  KEY `creator_user_id` (`creator_user_id`),
  KEY `assigned_to_user_id` (`assigned_to_user_id`),
  KEY `assigned_by_user_id` (`assigned_by_user_id`),
  KEY `work_ticket_ibfk_5` (`work_ticket_status_id`),
  CONSTRAINT `work_ticket_ibfk_1` FOREIGN KEY (`real_property_id`) REFERENCES `real_property` (`real_property_id`),
  CONSTRAINT `work_ticket_ibfk_2` FOREIGN KEY (`creator_user_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `work_ticket_ibfk_3` FOREIGN KEY (`assigned_to_user_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `work_ticket_ibfk_4` FOREIGN KEY (`assigned_by_user_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `work_ticket_ibfk_5` FOREIGN KEY (`work_ticket_type_id`) REFERENCES `work_ticket_type` (`work_ticket_type_id`),
  CONSTRAINT `work_ticket_ibfk_6` FOREIGN KEY (`work_ticket_status_id`) REFERENCES `work_ticket_status` (`work_ticket_status_id`),
  CONSTRAINT `work_ticket_ibfk_7` FOREIGN KEY (`ticket_group_id`) REFERENCES `ticket_group` (`ticket_group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `work_ticket`
--

LOCK TABLES `work_ticket` WRITE;
/*!40000 ALTER TABLE `work_ticket` DISABLE KEYS */;
INSERT INTO `work_ticket` VALUES (1,'75b2cc83091cd439d8f60832095567e',3,5,7,1,3,3,1,0,'2012-12-18 16:44:36','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'Title1','description1','2012-12-18 17:03:18'),(2,'8231704c868ef7ba24aa27463ba6522',4,1,7,1,4,4,1,0,'2012-12-18 17:50:18','1970-01-01 00:00:01',0.000000,0.000000,'1970-01-01 00:00:01',0.000000,0.000000,'eat lunch','beer','2012-12-18 17:50:18'),(3,'f8ba16a00160a03922dd01002d18733',3,5,7,1,3,3,1,0,'2012-12-18 18:24:15','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'title1','des1','2012-12-18 21:19:01'),(4,'6fd2c54d2240f7cb702383e8e9938fb',5,5,7,1,5,5,1,0,'2012-12-18 19:07:07','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'Work on Color Scheme','New Description','2012-12-18 21:40:43'),(5,'9bb5bc095be5b83824f96361343aaaa',3,5,7,1,3,3,1,0,'2012-12-18 19:41:59','2038-01-19 03:14:07',26.092289,-80.131683,'2012-12-23 15:55:42',26.092289,-80.131683,'title title title title title title title title title title title title title','fdsdsfsdf','2012-12-23 15:55:47'),(6,'55a66844dfa80ea9ba0137db2a238b8',3,5,7,1,3,3,1,0,'2012-12-18 19:44:07','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'cookie','','2012-12-18 21:21:17'),(7,'5d56e7ca0dafd6a83f7f33a8c2dcc6b',3,5,7,1,3,3,1,0,'2012-12-18 19:44:17','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'squibble','','2012-12-18 21:21:14'),(8,'e8177207402e7a390e02e3d15c7c0f6',3,5,7,1,3,3,1,0,'2012-12-18 19:44:37','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',0.000000,0.000000,'squable squable squable squable squable squable squable squable squable squable squable','','2012-12-18 21:40:58'),(9,'85ef00066aaa34dbea4e87c255f0a45',3,5,7,1,3,3,1,0,'2012-12-18 19:58:13','2038-01-19 03:14:07',26.092289,-80.131683,'2012-12-23 15:55:43',26.092289,-80.131683,'','','2012-12-23 15:55:47'),(10,'fade5ff80a4ae2691bf76512dbcc47c',3,5,7,1,3,3,1,0,'2012-12-18 21:21:30','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'test7','','2012-12-18 21:21:36'),(11,'93e26b219ae1005814afd4ecc783149',3,5,7,1,3,3,1,0,'2012-12-18 21:25:28','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'the title','','2012-12-18 21:25:42'),(13,'c1ea2a179a244a99c1c14907bcf8b04',5,5,7,1,5,5,1,0,'2012-12-18 21:37:48','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'New Ticket For Test','','2012-12-18 21:40:10'),(16,'85ae6344e7ff0219e50cbd2f8b9dace',3,2,7,1,3,3,1,0,'2012-12-18 21:39:10','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'cookie','','2012-12-18 21:39:17'),(17,'46b25c54ec152249dd102e0f026ee5f',3,5,7,1,3,3,1,0,'2012-12-19 13:43:09','2038-01-19 03:14:07',26.092289,-80.131683,'2038-01-19 03:14:07',26.092289,-80.131683,'title','fdsfdsfsd','2012-12-19 13:43:21'),(18,'d68b14d2190df888d7728dfebcad12d',7,5,7,1,7,7,1,0,'2012-12-19 16:52:47','1970-01-01 00:00:01',0.000000,0.000000,'2038-01-19 03:14:07',0.000000,0.000000,'Show off demo','Please don\'t crash','2012-12-19 16:53:35'),(19,'9c319ea083cc0ff929c86954e72fec2',7,5,7,1,7,7,1,0,'2012-12-19 16:53:27','1970-01-01 00:00:01',0.000000,0.000000,'2038-01-19 03:14:07',0.000000,0.000000,'','','2012-12-19 16:53:40'),(20,'b17c13605b3edb18cd35b96862c11c7',7,1,7,1,7,7,1,0,'2012-12-19 16:56:11','1970-01-01 00:00:01',0.000000,0.000000,'1970-01-01 00:00:01',0.000000,0.000000,'fdsdsffsd','','2012-12-19 16:56:11'),(21,'4acb2564005700287da7a197d660856',7,1,7,1,7,7,1,0,'2012-12-19 16:56:24','1970-01-01 00:00:01',0.000000,0.000000,'1970-01-01 00:00:01',0.000000,0.000000,'dis sad pdf fds. fds dgf fds fds pdf fds. fsdf','','2012-12-19 16:56:24'),(22,'ad08ba6156b15788b46c65aff5eb0b4',7,1,7,1,7,7,1,0,'2012-12-19 17:13:33','1970-01-01 00:00:01',0.000000,0.000000,'1970-01-01 00:00:01',0.000000,0.000000,'blah blah','do stuff','2012-12-19 17:13:33'),(23,'6c36cc047080ddc88a9ed097bfce13a',5,2,7,1,5,5,1,0,'2012-12-19 21:23:17','2038-01-19 03:14:07',26.092289,-80.131683,'1970-01-01 00:00:01',0.000000,0.000000,'Soemthing','','2012-12-19 21:23:24'),(24,'f4c02ef86017b10a872854e081e771b',3,5,7,1,3,3,1,0,'2012-12-23 16:01:51','2012-12-23 16:46:09',26.092289,-80.131683,'2012-12-23 16:46:28',26.092289,-80.131683,'Wire up invitation','connect invitation UI to API','2012-12-23 16:46:28'),(25,'409c1f2e5d3a53a943e7a39f90d439e',3,5,7,1,3,3,1,0,'2012-12-23 16:47:37','2012-12-23 19:06:49',26.092289,-80.131683,'2012-12-23 19:07:00',26.092289,-80.131683,'Setup current workspace','reconfigure top level workspace stuff','2012-12-23 19:07:01'),(26,'9aca84331fe0476bb8d93471129b75c',5,1,7,1,5,5,1,0,'2012-12-23 17:47:02','1970-01-01 00:00:01',0.000000,0.000000,'1970-01-01 00:00:01',0.000000,0.000000,'Something New','','2012-12-23 17:47:02'),(27,'9b43bc7f6cba63f85021dba6bfe32fc',5,1,7,1,5,5,1,0,'2012-12-23 20:30:21','1970-01-01 00:00:01',0.000000,0.000000,'1970-01-01 00:00:01',0.000000,0.000000,'ADdd New Ticket','','2012-12-23 20:30:21'),(28,'a1c10166d0bbdee9a5855eaeed63a65',5,5,7,1,5,5,1,0,'2012-12-23 20:30:38','2012-12-23 20:30:44',26.092289,-80.131683,'2012-12-23 20:30:49',26.092289,-80.131683,'vxcvxcvxcv','','2012-12-23 20:30:49');
/*!40000 ALTER TABLE `work_ticket` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `work_ticket_asset`
--

DROP TABLE IF EXISTS `work_ticket_asset`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `work_ticket_asset` (
  `work_ticket_asset_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `work_ticket_id` int(11) unsigned NOT NULL DEFAULT '1',
  `uploader_id` int(11) unsigned NOT NULL DEFAULT '1',
  `asset_name` varchar(32) NOT NULL DEFAULT '0',
  `content_type` varchar(32) NOT NULL DEFAULT '0',
  `asset_size_bytes` int(11) DEFAULT '0',
  `asset_size_height` int(11) DEFAULT '0',
  `asset_size_width` int(11) DEFAULT '0',
  `asset_upload_path` text,
  `asset_description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`work_ticket_asset_id`),
  KEY `work_ticket_id` (`work_ticket_id`),
  KEY `uploader_id` (`uploader_id`),
  KEY `asset_name` (`asset_name`),
  CONSTRAINT `work_ticket_asset_ibfk_1` FOREIGN KEY (`work_ticket_id`) REFERENCES `work_ticket` (`work_ticket_id`),
  CONSTRAINT `work_ticket_asset_ibfk_2` FOREIGN KEY (`uploader_id`) REFERENCES `end_user` (`end_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `work_ticket_asset`
--

LOCK TABLES `work_ticket_asset` WRITE;
/*!40000 ALTER TABLE `work_ticket_asset` DISABLE KEYS */;
INSERT INTO `work_ticket_asset` VALUES (1,1114,1,'1548868948999426642.png','image/png',123456,123456,123456,'','yadyaydyada','2012-07-19 15:55:15');
/*!40000 ALTER TABLE `work_ticket_asset` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `work_ticket_item`
--

DROP TABLE IF EXISTS `work_ticket_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `work_ticket_item` (
  `work_ticket_item_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `work_ticket_id` int(11) unsigned NOT NULL DEFAULT '1',
  `product_id` int(11) unsigned NOT NULL,
  `product_id_qty` float(10,4) NOT NULL,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`work_ticket_item_id`),
  KEY `product_id` (`product_id`),
  KEY `work_ticket_id` (`work_ticket_id`),
  CONSTRAINT `work_ticket_item_ibfk_1` FOREIGN KEY (`product_id`) REFERENCES `product` (`product_id`),
  CONSTRAINT `work_ticket_item_ibfk_2` FOREIGN KEY (`work_ticket_id`) REFERENCES `work_ticket` (`work_ticket_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `work_ticket_item`
--

LOCK TABLES `work_ticket_item` WRITE;
/*!40000 ALTER TABLE `work_ticket_item` DISABLE KEYS */;
INSERT INTO `work_ticket_item` VALUES (1,1,3,100.0000,'2012-12-23 20:28:20'),(2,1,3,100.0000,'2012-12-23 20:29:26'),(3,1,3,100.0000,'2012-12-23 20:30:50'),(4,1,3,100.0000,'2012-12-23 20:39:27'),(5,1,3,10.0000,'2012-12-23 20:46:00');
/*!40000 ALTER TABLE `work_ticket_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `work_ticket_status`
--

DROP TABLE IF EXISTS `work_ticket_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `work_ticket_status` (
  `work_ticket_status_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`work_ticket_status_id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `work_ticket_status`
--

LOCK TABLES `work_ticket_status` WRITE;
/*!40000 ALTER TABLE `work_ticket_status` DISABLE KEYS */;
INSERT INTO `work_ticket_status` VALUES (1,'OPEN','2012-08-16 15:49:34'),(2,'ACTIVE','2012-08-16 15:49:34'),(3,'PAUSED','2012-08-16 15:49:35'),(4,'BLOCKED','2012-08-16 15:51:27'),(5,'FINISHED','2012-08-16 15:51:27'),(6,'CLOSED','2012-08-16 15:51:27');
/*!40000 ALTER TABLE `work_ticket_status` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `work_ticket_type`
--

DROP TABLE IF EXISTS `work_ticket_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `work_ticket_type` (
  `work_ticket_type_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `description` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`work_ticket_type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `work_ticket_type`
--

LOCK TABLES `work_ticket_type` WRITE;
/*!40000 ALTER TABLE `work_ticket_type` DISABLE KEYS */;
INSERT INTO `work_ticket_type` VALUES (1,'general maintenance','2012-11-08 16:06:43'),(2,'plumbing','2012-11-08 16:06:43'),(3,'hvac','2012-11-08 16:06:43'),(4,'cleaning','2012-11-08 16:06:43'),(5,'painting','2012-11-08 16:06:43'),(6,'carpentry','2012-11-08 16:06:43'),(7,'electrical','2012-11-08 16:06:43'),(8,'software engineering','2012-05-16 18:24:26');
/*!40000 ALTER TABLE `work_ticket_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace`
--

DROP TABLE IF EXISTS `workspace`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `workspace` (
  `workspace_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_creator_id` int(11) unsigned NOT NULL DEFAULT '1',
  `workspace_owner_id` int(11) unsigned NOT NULL DEFAULT '1',
  `staging_group_id` int(11) unsigned NOT NULL DEFAULT '1',
  `workspace_token` varchar(32) NOT NULL DEFAULT '0',
  `workspace_name` varchar(64) NOT NULL DEFAULT '0',
  `workspace_description` text,
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`workspace_id`),
  UNIQUE KEY `workspace_token` (`workspace_token`),
  KEY `staging_group_id` (`staging_group_id`),
  KEY `workspace_creator_id` (`workspace_creator_id`),
  KEY `workspace_name` (`workspace_name`),
  KEY `workspace_ibfk_1` (`workspace_owner_id`),
  CONSTRAINT `workspace_ibfk_1` FOREIGN KEY (`workspace_owner_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `workspace_ibfk_2` FOREIGN KEY (`workspace_creator_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `workspace_ibfk_3` FOREIGN KEY (`staging_group_id`) REFERENCES `ticket_group` (`ticket_group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace`
--

LOCK TABLES `workspace` WRITE;
/*!40000 ALTER TABLE `workspace` DISABLE KEYS */;
INSERT INTO `workspace` VALUES (1,1,1,1,'eb170ea62a4448b4a609c0521fbb4cf9','Default','Default workspace for me','2012-07-09 17:26:08'),(2,2,2,2,'7830e3c1147541889ae357b4fcf8f3a0','Default','Default workspace for me','2012-05-07 18:39:11'),(4,3,3,3,'db2e5cec21b411ea95e10ece0304bc53','Workyplacey','Default Initial Workspace','2012-12-23 22:03:16'),(5,4,4,4,'b51c798321be11eaa2a812500a379b47','My First Workspace','Default Initial Workspace','2012-12-18 17:49:15'),(6,5,5,5,'964512dc21c511eab91612529ea83c2b','My First Workspace','Default Initial Workspace','2012-12-18 18:38:30'),(7,6,6,6,'bca033fd227e11ea9d3712cd07fb5023','My First Workspace','Default Initial Workspace','2012-12-19 16:43:51'),(8,7,7,7,'267e4198227f11ea9d3712cd07fb5023','My First Workspace','Default Initial Workspace','2012-12-19 16:46:49'),(9,8,8,8,'ec2fd77125a111ea90e706c4b7043921','My First Workspace','Default Initial Workspace','2012-12-23 16:33:17');
/*!40000 ALTER TABLE `workspace` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace_configuration`
--

DROP TABLE IF EXISTS `workspace_configuration`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `workspace_configuration` (
  `workspace_id` int(11) unsigned NOT NULL,
  `is_public` tinyint(1) NOT NULL DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`workspace_id`),
  CONSTRAINT `ws_configuration_ifbk_1` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`workspace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace_configuration`
--

LOCK TABLES `workspace_configuration` WRITE;
/*!40000 ALTER TABLE `workspace_configuration` DISABLE KEYS */;
/*!40000 ALTER TABLE `workspace_configuration` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace_members_lkp`
--

DROP TABLE IF EXISTS `workspace_members_lkp`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `workspace_members_lkp` (
  `wsmlkp_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `workspace_id` int(11) unsigned NOT NULL,
  `member_id` int(11) unsigned NOT NULL,
  `workspace_permission_id` int(11) unsigned NOT NULL,
  `active` tinyint(1) NOT NULL DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`wsmlkp_id`),
  UNIQUE KEY `grouping_comb1` (`workspace_id`,`member_id`),
  KEY `workspace_id` (`workspace_id`),
  KEY `member_id` (`member_id`),
  KEY `workspace_permission_id` (`workspace_permission_id`),
  KEY `last_updated` (`last_updated`),
  CONSTRAINT `ws_members_lkp_ibfk_1` FOREIGN KEY (`workspace_id`) REFERENCES `workspace` (`workspace_id`),
  CONSTRAINT `ws_members_lkp_ibfk_2` FOREIGN KEY (`member_id`) REFERENCES `end_user` (`end_user_id`),
  CONSTRAINT `ws_members_lkp_ibfk_3` FOREIGN KEY (`workspace_permission_id`) REFERENCES `workspace_permission` (`workspace_permission_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace_members_lkp`
--

LOCK TABLES `workspace_members_lkp` WRITE;
/*!40000 ALTER TABLE `workspace_members_lkp` DISABLE KEYS */;
INSERT INTO `workspace_members_lkp` VALUES (1,1,1,100,1,'2012-05-17 17:04:00'),(2,2,2,100,1,'2012-05-17 17:37:06'),(3,4,3,100,1,'2012-12-18 16:38:44'),(4,5,4,100,1,'2012-12-18 17:49:15'),(5,6,5,100,1,'2012-12-18 18:38:30'),(6,7,6,100,1,'2012-12-19 16:43:51'),(7,8,7,100,1,'2012-12-19 16:46:49'),(8,9,8,100,1,'2012-12-23 16:33:17'),(9,4,8,400,1,'2012-12-23 16:36:21');
/*!40000 ALTER TABLE `workspace_members_lkp` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `workspace_permission`
--

DROP TABLE IF EXISTS `workspace_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `workspace_permission` (
  `workspace_permission_id` int(11) unsigned NOT NULL,
  `workspace_permission` varchar(255) DEFAULT '0',
  `last_updated` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`workspace_permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `workspace_permission`
--

LOCK TABLES `workspace_permission` WRITE;
/*!40000 ALTER TABLE `workspace_permission` DISABLE KEYS */;
INSERT INTO `workspace_permission` VALUES (100,'ADMIN','2012-07-05 16:29:07'),(200,'DISPATCHER','2012-11-08 16:06:43'),(300,'SUPERVISOR','2012-11-08 16:06:43'),(400,'WORKER','2012-11-08 16:06:43');
/*!40000 ALTER TABLE `workspace_permission` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2012-10-31 16:10:39
